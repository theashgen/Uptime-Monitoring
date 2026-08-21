package checker

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/theashgen/url-short/internal/repo"
)

type CheckerService struct {
	queries *repo.Queries
}

type Job struct {
	Id uuid.UUID
	Url string
}


type Result struct {
	IsUp bool
	Error error
	StatusCode     int
	ResponseTimeMs int64
}

func NewCheckerService(queries *repo.Queries) *CheckerService {
	return &CheckerService{
		queries: queries,
	}
}

func Check(ctx context.Context, url string) Result {
	if  !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = "http://" + url
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{
			IsUp: false,
			Error: err,
		}
	}

	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		return Result{
			IsUp: false,
			Error: err,
			ResponseTimeMs: duration,
		}
	}
	defer resp.Body.Close()
	return Result{
		IsUp: true,
		StatusCode: resp.StatusCode,
		ResponseTimeMs: duration,
	}
}

func (s *CheckerService) Worker(ctx context.Context, jobs <- chan Job)  {
	for url := range jobs {
		res := Check(ctx, url.Url)
		
		var errString *string
		if res.Error != nil {
			errStr := res.Error.Error()
			errString = &errStr
		}
		s.queries.CreateURLCheck(ctx, repo.CreateURLCheckParams{
			UrlID:          url.Id,
			IsUp:           res.IsUp,
			Error:          errString,
			StatusCode:     res.StatusCode,
			ResponseTimeMs: res.ResponseTimeMs,
		})
		// fmt.Println("updated db")
		s.queries.UpdateURLNextCheck(ctx, url.Id)
	}
}



func (s *CheckerService) Scheduler(ctx context.Context) {
	
	n_worker := 5
	jobs := make(chan Job)

	for i := 0; i < n_worker; i ++ {
		go s.Worker(ctx, jobs)
	}

	for {
		// fmt.Println("urls")
		urls, err := s.queries.GetDueURLs(ctx, 100)
		// fmt.Print(urls)
		if err != nil {
			// fmt.Println("Error while getting due urls.")
			time.Sleep(time.Second * 5)
			continue	
		}
		
		for _, url := range urls {
			jobs <- Job {
				Id: url.ID,
				Url: url.Url,
			}
		}

		// Wait for the 5s
		// fmt.Println("waiting for next check")
		time.Sleep(time.Minute)
	
	}	
}


