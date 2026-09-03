package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/theashgen/url-short/internal/repo"
)

// urlChecksHistoryLimit caps how many recent check results are returned for a
// single URL's status. Checks are best-effort telemetry (see checker package),
// so there's no pagination — just the most recent window.
const urlChecksHistoryLimit = 50

// URLStatus is a URL's metadata plus its recent check history. It's an
// internal service/handler boundary type, not serialized directly — the
// handler maps it into its own response contract (handler.GetUrlStatusResponse).
type URLStatus struct {
	URL    repo.Url
	Checks []repo.UrlCheck
}

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

func (s *URLService) GetURLStatusByUsername(ctx context.Context, username string, urlID string) (URLStatus, error) {
	id, err := uuid.Parse(urlID)
	if err != nil {
		return URLStatus{}, fmt.Errorf("%w: invalid url id", ErrInvalidInput)
	}

	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return URLStatus{}, fmt.Errorf("%w: user", ErrNotFound)
		}
		return URLStatus{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	// Scoped by user_id, not just id, so one user can't probe another's URL by
	// guessing/enumerating UUIDs — a mismatch looks identical to a missing URL.
	url, err := s.queries.GetURLByID(ctx, repo.GetURLByIDParams{ID: id, UserID: user.ID})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return URLStatus{}, fmt.Errorf("%w: url", ErrNotFound)
		}
		return URLStatus{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	checks, err := s.queries.ListURLChecksByURL(ctx, repo.ListURLChecksByURLParams{
		UrlID: id,
		Limit: urlChecksHistoryLimit,
	})
	if err != nil {
		return URLStatus{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return URLStatus{URL: url, Checks: checks}, nil
}
