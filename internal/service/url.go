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
