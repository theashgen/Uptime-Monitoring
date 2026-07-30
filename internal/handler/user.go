package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/theashgen/url-short/internal/auth"
	"github.com/theashgen/url-short/internal/service"
)

type UserHandler struct {
	userHandler *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userHandler: userService,
	}
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignRespones struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *UserHandler) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var userBody UserSignUpRequest

	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		http.Error(w, "Invalid json request", http.StatusBadRequest)
		return
	}

	user, err := h.userHandler.CreateUser(r.Context(), userBody.Email, userBody.Username, userBody.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}



func (h *UserHandler) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var UserLoginDetails UserLoginRequest

	err := json.NewDecoder(r.Body).Decode(&UserLoginDetails)
	if err != nil {
		http.Error(w, "Invalid json format", http.StatusBadRequest)
		return
	}
	user, err := h.userHandler.AuthenticateUser(r.Context(), UserLoginDetails.Email, UserLoginDetails.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := auth.SignUserJWT(user.Username, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(w, "error while creating token", http.StatusInternalServerError)
		return
	}

	auth.SetCookie(w, token)
	
	json.NewEncoder(w).Encode(map[string]string{
		"message": "login successful",
	})

}
