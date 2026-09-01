package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

type InsertURLbyUsernameParams struct {
	Username string
	Url      string
	Interval int32
}

func (s *URLService) ListURLsByUsername(ctx context.Context, username string) ([]repo.ListURLsByUserRow, error) {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: user", ErrNotFound)
		}
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// A user with no monitored URLs yet is a valid, non-error outcome (an empty list),
	// not a "not found" condition — sqlc :many queries return an empty slice + nil error
	// in that case, so there's nothing further to check here.
	urls, err := s.queries.ListURLsByUser(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return urls, nil
}

func (s *URLService) InsertURLbyUsername(ctx context.Context, params InsertURLbyUsernameParams) (repo.CreateURLRow, error) {
	if params.Username == "" {
		return repo.CreateURLRow{}, fmt.Errorf("%w: provide a username", ErrInvalidInput)
	}

	user, err := s.queries.GetUserByUsername(ctx, params.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CreateURLRow{}, fmt.Errorf("%w: user", ErrNotFound)
		}
		return repo.CreateURLRow{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	url, err := s.queries.CreateURL(ctx, repo.CreateURLParams{
		UserID:          user.ID,
		Url:             params.Url,
		IntervalSeconds: params.Interval,
	})
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return repo.CreateURLRow{}, fmt.Errorf("%w: host already being monitored", ErrConflict)
		}
		return repo.CreateURLRow{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}
	return url, nil
}
