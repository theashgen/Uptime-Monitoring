package service

import (
	"context"
	"errors"

	"github.com/theashgen/url-short/internal/repo"
)

type URLService struct {
	queries *repo.Queries
}

func NewURLService(queries *repo.Queries) *URLService {
	return &URLService{
		queries: queries,
	}
}

type InsertURLbyUsernameParams struct{
	Username string
	Host string
	Interval string
}

func (s *URLService) ListURLsByUsername(ctx context.Context, username string) ([]repo.ListURLsByUserRow, error) {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("User doesnt exists")
	}

	urls, err := s.queries.ListURLsByUser(ctx, user.ID)
	if err != nil {
		return nil, errors.New("No urls has been yet created")
	}

	return urls, nil
}


func (s *URLService) InsertURLbyUsername(ctx context.Context, params InsertURLbyUsernameParams) (repo.CreateURLRow, error) {
	if params.Username == "" {
		return repo.CreateURLRow{}, errors.New("provide username")
	}

	user, err := s.queries.GetUserByUsername(ctx, params.Username)
	if err != nil {
		return repo.CreateURLRow{}, err
	}
	url, err := s.queries.CreateURL(ctx, repo.CreateURLParams{
		UserID: user.ID,
		Host: params.Host,
		Interval: params.Interval,
	})
	
	if err != nil {
		return repo.CreateURLRow{}, err
	}

	return url, nil
}
