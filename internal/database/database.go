package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://postgres:postgres@localhost:5432/urlshortener",
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
