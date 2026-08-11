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

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (repo.GetUserByUsernameRow, error) {
	if username == "" {
		return repo.GetUserByUsernameRow{}, errors.New("provide valid username.")
	}

	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return repo.GetUserByUsernameRow{}, errors.New("username doesnt exist in the database")
	}

	return user, nil
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
) (repo.GetUserByEmailRow, error) {
	// start := time.Now()
	user, err := s.queries.GetUserByEmail(ctx, email)

	// fmt.Println("DB:", time.Since(start))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repo.GetUserByEmailRow{}, errors.New("invalid email or password")
		}
		return repo.GetUserByEmailRow{}, err
	}

	// start = time.Now()

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Passwordhash),
		[]byte(password),
	); err != nil {
		return repo.GetUserByEmailRow{}, errors.New("invalid email or password")
	}
	// fmt.Println("bcrypt:", time.Since(start))

	return user, nil
}
