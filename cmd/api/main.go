package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/theashgen/url-short/internal/database"
	"github.com/theashgen/url-short/internal/handler"
	"github.com/theashgen/url-short/internal/middleware"
	"github.com/theashgen/url-short/internal/repo"
	"github.com/theashgen/url-short/internal/service"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(context.Background())

	queries := repo.New(db)

	userService := service.NewUserService(queries)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/signup", userHandler.UserSignUp)
	mux.HandleFunc("POST /api/v1/login", userHandler.UserLoginHandler)

	handler := middleware.Logger(mux)

	log.Println("server running on :3000")

	err = http.ListenAndServe(":3000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
