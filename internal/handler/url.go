package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/theashgen/url-short/internal/repo"
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
	Url string `json:"url"`
	Interval int32 `json:"interval"`
}

// GetUrls godoc
// @Summary      Get URLs
// @Description  Get all monitored URLs for the authenticated user
// @Tags         urls
// @Produce      json
// @Security     CookieAuth
// @Success      200  {array}   repo.ListURLsByUserRow
// @Failure      401  {string}  string "username not found"
// @Failure      404  {string}  string "url not found"
// @Failure      500  {string}  string "internal server error"
// @Router       /urls [get]
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

// PostUrl godoc
// @Summary      Post URL
// @Description  Add a new URL to monitor for the authenticated user
// @Tags         urls
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        request body handler.PostUrlBody true "Add URL Request Body"
// @Success      200  {object}  repo.CreateURLRow
// @Failure      401  {string}  string "Unauthorized"
// @Failure      400  {string}  string "Invalid json input / interval should be greater than zero."
// @Failure      500  {string}  string "internal server error"
// @Router       /urls [post]
func (h *URLHandler) PostUrl(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
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

	if reqBody.Interval <= 0 {
		http.Error(w, "interval should be greater than zero.", http.StatusBadRequest)
		return
	}
	url, err := h.urlService.InsertURLbyUsername(r.Context(), service.InsertURLbyUsernameParams{
		Username: username,
		Url: reqBody.Url,
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
