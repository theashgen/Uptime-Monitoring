package handler

import (
	"github.com/theashgen/url-short/internal/service"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}
