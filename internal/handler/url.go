package handler

import (
	"encoding/json"
	"fmt"
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

type PostUrlBody struct {
	Host string `json:"host"`
	Interval string `json:"interval"`
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

	w.Header().Set("Content-Type", "application/json")	
	err = json.NewEncoder(w).Encode(urls)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")	
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *URLHandler) PostUrl(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	fmt.Println("username:", username)
	fmt.Println("ok:", ok)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	var reqBody PostUrlBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid json input", http.StatusBadRequest)
		return 
	}

	url, err := h.urlService.InsertURLbyUsername(r.Context(), service.InsertURLbyUsernameParams{
		Username: username,
		Host: reqBody.Host,
		Interval: reqBody.Interval,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(url)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
