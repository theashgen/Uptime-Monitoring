package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/theashgen/url-short/internal/helper"
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
	Username string `json:"username"`
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

	token, ss = helper.SignUserJWT()
	fmt.Fprint(w, "need to code this route")

	
}
