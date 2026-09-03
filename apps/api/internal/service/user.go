package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/theashgen/url-short/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	queries *repo.Queries
}

func NewUserService(queries *repo.Queries) *UserService {
	return &UserService{
		queries: queries,
	}
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (repo.GetUserByUsernameRow, error) {
	if username == "" {
		return repo.GetUserByUsernameRow{}, fmt.Errorf("%w: provide a valid username", ErrInvalidInput)
	}

	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.GetUserByUsernameRow{}, fmt.Errorf("%w: user", ErrNotFound)
		}
		return repo.GetUserByUsernameRow{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, email, username, password string) (repo.CreateUserRow, error) {

	if username == "" || password == "" || email == "" {
		return repo.CreateUserRow{}, fmt.Errorf("%w: provide a valid username, password and email", ErrInvalidInput)
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return repo.CreateUserRow{}, fmt.Errorf("%w: hashing password: %v", ErrInternal, err)
	}

	user, err := s.queries.CreateUser(ctx, repo.CreateUserParams{
		Username:     username,
		Email:        email,
		Passwordhash: string(passwordHash),
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return repo.CreateUserRow{}, fmt.Errorf("%w: user already exists", ErrConflict)
		}
		return repo.CreateUserRow{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(
	ctx context.Context,
	email, password string,
) (repo.GetUserByEmailRow, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.GetUserByEmailRow{}, ErrUnauthenticated
		}
		return repo.GetUserByEmailRow{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Passwordhash),
		[]byte(password),
	); err != nil {
		return repo.GetUserByEmailRow{}, ErrUnauthenticated
	}

	return user, nil
}
