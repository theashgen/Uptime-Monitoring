//	@title			Uptime Monitor API
//	@version		1.0
//	@description	API Server for monitoring website uptime.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/api/v1

// @securityDefinitions.apikey	CookieAuth
// @in							cookie
// @name						access_token
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/theashgen/url-short/docs"
	"github.com/theashgen/url-short/internal/database"
	"github.com/theashgen/url-short/internal/handler"
	"github.com/theashgen/url-short/internal/middleware"
	"github.com/theashgen/url-short/internal/repo"
	"github.com/theashgen/url-short/internal/service"
	"github.com/theashgen/url-short/internal/service/checker"
)

func main() {

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	queries := repo.New(db)

	ctx := context.Background()
	s := checker.NewCheckerService(queries)
	go s.Scheduler(ctx) // <- if !routine block thread

	userService := service.NewUserService(queries)
	userHandler := handler.NewUserHandler(userService)

	urlService := service.NewURLService(queries)
	urlHandler := handler.NewURLHandler(urlService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/signup", userHandler.UserSignUp)
	mux.HandleFunc("POST /api/v1/login", userHandler.UserLoginHandler)

	handler := middleware.Logger(mux)

	mux.Handle("GET /api/v1/urls",
		middleware.AuthMiddleware(urlHandler.GetUrls),
	)

	mux.Handle("POST /api/v1/urls",
		middleware.AuthMiddleware(urlHandler.PostUrl),
	)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	log.Println("server running on :3000")

	err = http.ListenAndServe(":3000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
