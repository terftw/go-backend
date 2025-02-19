package main

import (
	// "context"
	"fmt"
	"log"
	"net/http"

	// "net/http"
	// "os"
	// "os/signal"
	// "time"

	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	// "github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/terftw/go-backend/internal/api/handlers"
	"github.com/terftw/go-backend/internal/api/routes"
	"github.com/terftw/go-backend/internal/config"
	"github.com/terftw/go-backend/internal/db"
	"github.com/terftw/go-backend/internal/db/repositories"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db.Connect(&config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repos := repositories.NewRepository(db)
	handlers := handlers.NewHandlers(repos.UserRepository, config.OAuth.GoogleOAuth, config.PrivateKey)

	router := chi.NewRouter()
	routes.SetupRoutes(router, handlers)

	addr := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Server starting on port %d", config.Server.Port)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
