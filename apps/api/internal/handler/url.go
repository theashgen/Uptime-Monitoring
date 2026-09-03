package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/theashgen/url-short/internal/middleware"
	"github.com/theashgen/url-short/internal/repo"
	"github.com/theashgen/url-short/internal/service"
)

// writeServiceError maps a service-layer sentinel error to the appropriate HTTP status.
// Internal errors are logged with their real detail but never exposed to the client.
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		http.Error(w, "invalid input", http.StatusBadRequest)
	case errors.Is(err, service.ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, service.ErrConflict):
		http.Error(w, "already exists", http.StatusConflict)
	case errors.Is(err, service.ErrUnauthenticated):
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

type PostUrlRequest struct {
	Url      string `json:"url"`
	Interval int32  `json:"interval"`
}

// UrlResponse is a monitored URL's public metadata. It's a handler-owned
// contract, deliberately separate from repo.Url/repo.ListURLsByUserRow, so a
// schema/query change doesn't automatically change the API response shape.
type UrlResponse struct {
	ID              uuid.UUID `json:"id"`
	Url             string    `json:"url"`
	IntervalSeconds int32     `json:"intervalSeconds"`
	NextCheckAt     time.Time `json:"nextCheckAt"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
}

func toUrlResponse(u repo.ListURLsByUserRow) UrlResponse {
	return UrlResponse{
		ID:              u.ID,
		Url:             u.Url,
		IntervalSeconds: u.IntervalSeconds,
		NextCheckAt:     u.NextCheckAt.Time,
		IsActive:        u.IsActive,
		CreatedAt:       u.CreatedAt.Time,
	}
}

type PostUrlResponse struct {
	ID              uuid.UUID `json:"id"`
	Url             string    `json:"url"`
	IntervalSeconds int32     `json:"intervalSeconds"`
}

type UrlCheckResponse struct {
	ID             uuid.UUID `json:"id"`
	IsUp           bool      `json:"isUp"`
	StatusCode     int       `json:"statusCode"`
	ResponseTimeMs int64     `json:"responseTimeMs"`
	Error          *string   `json:"error"`
	CheckedAt      time.Time `json:"checkedAt"`
}

func toUrlCheckResponse(c repo.UrlCheck) UrlCheckResponse {
	return UrlCheckResponse{
		ID:             c.ID,
		IsUp:           c.IsUp,
		StatusCode:     c.StatusCode,
		ResponseTimeMs: c.ResponseTimeMs,
		Error:          c.Error,
		CheckedAt:      c.CheckedAt.Time,
	}
}

type GetUrlStatusResponse struct {
	Url    UrlResponse        `json:"url"`
	Checks []UrlCheckResponse `json:"checks"`
}

// GetUrls godoc
//
//	@Summary		Get URLs
//	@Description	Get all monitored URLs for the authenticated user
//	@Tags			urls
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{array}		handler.UrlResponse
//	@Failure		401	{string}	string	"username not found"
//	@Failure		404	{string}	string	"url not found"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/urls [get]
func (h *URLHandler) GetUrls(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "username not found", http.StatusUnauthorized)
		return
	}

	urls, err := h.urlService.ListURLsByUsername(r.Context(), username)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	resp := make([]UrlResponse, len(urls))
	for i, u := range urls {
		resp[i] = toUrlResponse(u)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// PostUrl godoc
//
//	@Summary		Post URL
//	@Description	Add a new URL to monitor for the authenticated user
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			request	body		handler.PostUrlRequest	true	"Add URL Request Body"
//	@Success		200		{object}	handler.PostUrlResponse
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		400		{string}	string	"Invalid json input / interval should be greater than zero."
//	@Failure		500		{string}	string	"internal server error"
//	@Router			/urls [post]
func (h *URLHandler) PostUrl(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody PostUrlRequest
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
		Url:      reqBody.Url,
		Interval: reqBody.Interval,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	resp := PostUrlResponse{
		ID:              url.ID,
		Url:             url.Url,
		IntervalSeconds: url.IntervalSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// GetUrlStatus godoc
//
//	@Summary		Get URL status
//	@Description	Get a monitored URL's metadata and recent check history for the authenticated user
//	@Tags			urls
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"URL ID"
//	@Success		200	{object}	handler.GetUrlStatusResponse
//	@Failure		400	{string}	string	"invalid url id"
//	@Failure		401	{string}	string	"username not found"
//	@Failure		404	{string}	string	"url not found"
//	@Failure		500	{string}	string	"internal server error"
//	@Router			/urls/{id} [get]
func (h *URLHandler) GetUrlStatus(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UsernameKey).(string)
	if !ok {
		http.Error(w, "username not found", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")

	status, err := h.urlService.GetURLStatusByUsername(r.Context(), username, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	checks := make([]UrlCheckResponse, len(status.Checks))
	for i, c := range status.Checks {
		checks[i] = toUrlCheckResponse(c)
	}

	resp := GetUrlStatusResponse{
		Url: UrlResponse{
			ID:              status.URL.ID,
			Url:             status.URL.Url,
			IntervalSeconds: status.URL.IntervalSeconds,
			NextCheckAt:     status.URL.NextCheckAt.Time,
			IsActive:        status.URL.IsActive,
			CreatedAt:       status.URL.CreatedAt.Time,
		},
		Checks: checks,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
