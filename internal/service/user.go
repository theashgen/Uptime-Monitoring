package service

import (
	"context"
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

func (s *UserService) GetUserFromEmail(ctx context.Context, email, password string) (repo.GetUserEmailRow, error) {
	user, err := s.queries.GetUserEmail(ctx, email)
	if err != nil {
		return repo.GetUserEmailRow{}, errors.New("Internal Server Err")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Passwordhash), []byte(password))
	if err != nil {
		return repo.GetUserEmailRow{}, errors.New("Password is wrong")
	}
	return user, nil
}
