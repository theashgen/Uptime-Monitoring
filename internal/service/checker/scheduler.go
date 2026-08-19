package checker

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"strings"

	"github.com/theashgen/url-short/internal/repo"
)

type CheckerService struct {
	queries *repo.Queries
}

func NewCheckerService(queries *repo.Queries) *CheckerService {
	return &CheckerService{
		queries: queries,
	}
}

type Result struct {
	IsUp bool
	Error error
	StatusCode     int
	ResponseTimeMs int64
}

func Check(ctx context.Context, url string) Result {

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
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
	responseTime := time.Since(start).Milliseconds()
	if err != nil {
		return Result{
			IsUp: false,
			Error: err,
			ResponseTimeMs: responseTime,
		}
	}
	defer resp.Body.Close()
	
	return Result{
		IsUp: resp.StatusCode >= 200 && resp.StatusCode < 400,
		Error: nil,
		StatusCode: resp.StatusCode,
		ResponseTimeMs: responseTime,
	}
}

func (s *CheckerService) Scheduler(ctx context.Context) {
	for {

		urls, err := s.queries.GetDueURLs(ctx, 10)
		if err != nil {
			fmt.Println("Error while getting due urls.")
			time.Sleep(time.Second * 5)
			continue	
		}
		
		for _, url := range urls {
			res := Check(ctx, url.Url)	
			
			var errString *string
			if res.Error != nil {
				errstr := res.Error.Error()
				errString = &errstr
			}

			_, err = s.queries.CreateURLCheck(ctx, repo.CreateURLCheckParams{
				UrlID: url.ID,
				IsUp: res.IsUp,
				Error: errString,
				StatusCode: res.StatusCode,
				ResponseTimeMs: res.ResponseTimeMs,
			})
			if err != nil {
				fmt.Println(err)
			}

			err = s.queries.UpdateURLNextCheck(ctx, url.ID)
			if err != nil {
				fmt.Print(err.Error())
			}
		}

		// Wait for the 5s
		fmt.Println("waiting for next check")
		time.Sleep(time.Second * 5)
	
	}	
}


