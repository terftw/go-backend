package main

import (
	// "context"
	// "fmt"
	"log"
	// "net/http"
	// "os"
	// "os/signal"
	// "time"

	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	// "github.com/go-chi/cors"
	"github.com/joho/godotenv"
	// "gorm.io/gorm"

	"github.com/terftw/go-backend/internal/config"
	"github.com/terftw/go-backend/internal/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db.NewConnection(&config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var result int
	db.Raw("SELECT 1").Scan(&result)
	log.Printf("Database connection test: %d", result)

	go func() {
		log.Printf("Server starting on port %d", config.Server.Port)
	}()
}
