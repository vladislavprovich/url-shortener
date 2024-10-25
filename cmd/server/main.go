package main

import (
	"context"
	"database/sql"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vladislavprovich/url-shortener/internal/handler"
	"github.com/vladislavprovich/url-shortener/internal/middleware"
	"github.com/vladislavprovich/url-shortener/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// TODO move to separate func
	// Initialize Logger
	logger := initLoger()
	db := initDB()
	defer func() {

		if err := db.Close(); err != nil {
			logger.Warn("Error closing db", zap.Error(err))
		}

	}()

	//todo fix error
	repo := initRepo(db)
	service := initService(&repo, logger)
	urlHandler := initHandler(service, logger)

	// TODO move to separate func
	// Connect to Database
	// Create Router
	r := chi.NewRouter()

	// Apply Middlewares
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter)

	// TODO repository and service must init here together with urlHandler
	// Initialize Handlers

	// Routes
	r.Post("/shorten", urlHandler.ShortenURL)
	r.Get("/{shortURL}", urlHandler.Redirect)
	r.Get("/{shortURL}/stats", urlHandler.GetStats)

	// Start Server Graceful Shutdown
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Server is starting", zap.String("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	// Block until a signal is received
	<-quit
	logger.Info("Server is shutting down...")

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

func initLoger() *zap.Logger {
	// Initialize Logger
	return logger.NewLogger(os.Getenv("LOG_LEVEL"))

}

func initDB() *sql.DB {
	//Connect to Database
	db, err := repository.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	return db
}

func initRepo(db *sql.DB) repository.URLRepository {
	return repository.NewURLRepository(db)
}

func initService(repo *repository.URLRepository, logger *zap.Logger) service.URLService {
	return service.NewURLService(*repo, logger)

}

func initHandler(srv service.URLService, logger *zap.Logger) *handler.URLHandler {
	return handler.NewURLHandler(srv, logger)
}
