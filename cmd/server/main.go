package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/vladislavprovich/url-shortener/internal/handler"
	"github.com/vladislavprovich/url-shortener/internal/middleware"
	"github.com/vladislavprovich/url-shortener/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// TODO move to separate func
	// Initialize Logger
	logger := logger.NewLogger(os.Getenv("LOG_LEVEL"))

	// TODO move to separate func
	// Connect to Database
	db, err := handler.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Create Router
	r := chi.NewRouter()

	// Apply Middlewares
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter)

	// Initialize Handlers
	urlHandler := handler.NewURLHandler(db, logger)

	// Routes
	r.Post("/shorten", urlHandler.ShortenURL)
	r.Get("/{shortURL}", urlHandler.Redirect)
	r.Get("/{shortURL}/stats", urlHandler.GetStats)

	// Start Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
