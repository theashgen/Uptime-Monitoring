package checker

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theashgen/url-short/internal/repo"
)

type CheckerService struct {
	queries *repo.Queries
	client  *http.Client
}

type Job struct {
	Id  uuid.UUID
	Url string
}

type CheckResult struct {
	URLID          uuid.UUID
	IsUp           bool
	Error          *string
	StatusCode     int32
	ResponseTimeMs int64
}

type Result struct {
	IsUp           bool
	Error          error
	StatusCode     int
	ResponseTimeMs int64
}

func NewCheckerService(queries *repo.Queries) *CheckerService {
	return &CheckerService{
		queries: queries,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *CheckerService) Check(ctx context.Context, url string) Result {
	if !(strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://")) {
		url = "http://" + url
	}

	start := time.Now()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return Result{
			IsUp:  false,
			Error: err,
		}
	}

	resp, err := s.client.Do(req)

	duration := time.Since(start).Milliseconds()

	if err != nil {
		return Result{
			IsUp:           false,
			Error:          err,
			ResponseTimeMs: duration,
		}
	}

	defer resp.Body.Close()

	return Result{
		IsUp:           resp.StatusCode >= 200 && resp.StatusCode < 400,
		StatusCode:     resp.StatusCode,
		ResponseTimeMs: duration,
	}
}

func (s *CheckerService) Worker(
	ctx context.Context,
	jobs <-chan Job,
	results chan<- repo.CreateURLChecksParams,
) {
	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			res := s.Check(ctx, job.Url)

			var errString *string
			if res.Error != nil {
				err := res.Error.Error()
				errString = &err
			}

			results <- repo.CreateURLChecksParams{
				UrlID:          job.Id,
				IsUp:           res.IsUp,
				Error:          errString,
				StatusCode:     res.StatusCode,
				ResponseTimeMs: res.ResponseTimeMs,
			}
		}
	}
}

func (s *CheckerService) ResultProcessor(
	ctx context.Context,
	results <-chan repo.CreateURLChecksParams,
) {
	batch := make([]repo.CreateURLChecksParams, 0, 100)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	flush := func(flushCtx context.Context) {
		if len(batch) == 0 {
			return
		}

		_, err := s.queries.CreateURLChecks(flushCtx, batch)
		if err != nil {
			log.Printf("bulk insert failed: %v", err)
		}

		// Drop the batch either way: check results are best-effort telemetry, not critical
		// data, so we bound memory instead of retrying/growing the batch forever on a
		// persistent DB outage.
		batch = batch[:0]
	}

	for {
		select {
		case <-ctx.Done():
			// ctx is already cancelled here, so the final flush needs its own short-lived
			// context instead of the cancelled one (which would fail immediately).
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			flush(flushCtx)
			cancel()
			return

		case result := <-results:
			batch = append(batch, result)

			if len(batch) >= 100 {
				flush(ctx)
			}

		case <-ticker.C:
			flush(ctx)
		}
	}
}


func (s *CheckerService) Scheduler(ctx context.Context) {
	
	n_worker := 100
	jobs := make(chan Job, 100)
	defer close(jobs)

	results := make(chan repo.CreateURLChecksParams, 100)

	for i := 0; i < n_worker; i ++ {
		go s.Worker(ctx, jobs, results)
	}

	go s.ResultProcessor(ctx, results)

	for {
		// fmt.Println("urls")
		urls, err := s.queries.ClaimDueURLs(ctx, 100)
		// fmt.Print(urls)
		if err != nil {
			// fmt.Println("Error while getting due urls.")
			time.Sleep(time.Second * 5)
			continue
		}

		if len(urls) == 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}

		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- Job{
					Id:  url.ID,
					Url: url.Url,
				}:
			}
		}
	}	
}
