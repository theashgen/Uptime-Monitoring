package checker

import (
	"context"
	"net/http"
	"time"

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
		
			return
		}
		
		for _, url := range urls {
			res := Check(ctx, url.Url)	
			
			s.queries.CreateURLCheck(ctx, repo.CreateURLCheckParams{
				UrlID: url.ID,
				IsUp: res.IsUp,
				StatusCode: res.StatusCode,
			})
		}

		// Wait for the 5s
		time.Sleep(time.Second * 5)
	
	}	
}


