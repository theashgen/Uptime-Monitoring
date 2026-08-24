package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/theashgen/url-short/internal/auth"
	_ "github.com/theashgen/url-short/internal/repo"
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

// UserSignUp godoc
// @Summary      User SignUp
// @Description  Create a new user account
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body handler.UserSignUpRequest true "Sign Up Request Body"
// @Success      200  {object}  repo.CreateUserRow
// @Failure      400  {string}  string "Invalid json request"
// @Failure      409  {string}  string "User already exists"
// @Failure      500  {string}  string "Internal Server error"
// @Router       /signup [post]
func (h *UserHandler) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var userBody UserSignUpRequest

	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		http.Error(w, "Invalid json request", http.StatusBadRequest)
		return
	}

	user, err := h.userHandler.CreateUser(r.Context(), userBody.Email, userBody.Username, userBody.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")	
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")	
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}

// UserLoginHandler godoc
// @Summary      User Login
// @Description  Authenticate user and set cookie
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body handler.UserLoginRequest true "Login Request Body"
// @Success      200  {object}  map[string]string "{"message": "login successful"}"
// @Failure      400  {string}  string "Invalid json format"
// @Failure      401  {string}  string "invalid email or password"
// @Failure      500  {string}  string "error while creating token / json parsing failed"
// @Router       /login [post]
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
	w.Header().Set("Content-Type", "application/json")	
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "login successful",
	})
	
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")	
		http.Error(w, "json parasing failed", http.StatusInternalServerError)
	}
}
