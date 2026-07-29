package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/theashgen/url-short/internal/database"
	"github.com/theashgen/url-short/internal/handler"
	"github.com/theashgen/url-short/internal/repo"
	"github.com/theashgen/url-short/internal/service"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	queries := repo.New(db)

	// urlService := service.NewURLService(queries)
	userService := service.NewUserService(queries)

	userHandler := handler.NewUserHandler(userService)
	// urlHandler := handler.NewURLHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/signup", userHandler.UserSignUp)
	mux.HandleFunc("POST /api/v1/login", userHandler.UserLoginHandler)

	defer db.Close(context.Background())

	http.ListenAndServe(":3000", mux)
}
