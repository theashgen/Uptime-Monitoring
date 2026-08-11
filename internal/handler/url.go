package handler

import (
	"encoding/json"
	"net/http"

	"github.com/theashgen/url-short/internal/middleware"
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

func (h *URLHandler) GetUrls(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "username not found", http.StatusUnauthorized)
		return
	}

	urls, err := h.urlService.ListURLsByUsername(r.Context(), username)
	if err != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(urls)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) PostUrl(w http.ResponseWriter, r *http.Request) {
	
}
