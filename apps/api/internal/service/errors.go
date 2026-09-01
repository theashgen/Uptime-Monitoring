package service

import "errors"

// Sentinel errors returned by this package. Handlers should check these with errors.Is
// to decide the HTTP status code, instead of pattern-matching error strings or reaching
// into DB-driver error types (pgconn.PgError, pgx.ErrNoRows, etc.) themselves.
var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("already exists")
	ErrUnauthenticated = errors.New("invalid email or password")
	ErrInternal        = errors.New("internal error")
)
