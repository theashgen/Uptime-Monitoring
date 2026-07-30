package service

import (
	"context"
	"database/sql"
	"errors"

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

func (s *UserService) CreateUser(ctx context.Context, email, username, password string) (repo.CreateUserRow, error) {

	if username == "" || password == "" || email == "" {
		return repo.CreateUserRow{}, errors.New("provide valid username, password and gmail")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return repo.CreateUserRow{}, errors.New("Interal hashing failed")
	}

	user, err := s.queries.CreateUser(ctx, repo.CreateUserParams{
		Username:     username,
		Email:        email,
		Passwordhash: string(passwordHash),
	})
	if err != nil {
		return repo.CreateUserRow{}, errors.New("User already exisits")
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(
	ctx context.Context,
	email, password string,
) (repo.GetUserEmailRow, error) {
	// start := time.Now()
	user, err := s.queries.GetUserEmail(ctx, email)

	// fmt.Println("DB:", time.Since(start))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.GetUserEmailRow{}, errors.New("invalid email or password")
		}
		return repo.GetUserEmailRow{}, err
	}

	// start = time.Now()

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Passwordhash),
		[]byte(password),
	); err != nil {
		return repo.GetUserEmailRow{}, errors.New("invalid email or password")
	}
	// fmt.Println("bcrypt:", time.Since(start))

	return user, nil
}
