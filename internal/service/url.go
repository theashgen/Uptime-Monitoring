package service

import (
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
