package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		context.Background(),
		"postgres://ashwin@localhost:5432/urlshort?sslmode=disable",
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
