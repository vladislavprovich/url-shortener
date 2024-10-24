package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/internal/service"
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

var log *zap.Logger
var db *sql.DB
var repo *repository.URLRepository

func main() {
	// TODO move to separate func
	// Initialize Logger
	// TODO move to separate func
	// Connect to Database
	// Create Router
	r := chi.NewRouter()

	// Apply Middlewares
	r.Use(middleware.Recoverer(log))
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter)

	// TODO repository and service must init here together with urlHandler
	// Initialize Handlers
	urlHandler := handler.NewURLHandler(db, log)
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
		log.Info("Server is starting", zap.String("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error", zap.Error(err))
		}
	}()

	// Block until a signal is received
	<-quit
	log.Info("Server is shutting down...")

	// Create a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}

func initLoger() {
	// Initialize Logger
	log := logger.NewLogger(os.Getenv("LOG_LEVEL"))
	if log == nil {
		log.Fatal("Failed to initialize logger")
	}

}

func initDataBase() {
	// Connect to Database
	db, err := handler.InitDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()
}

func initRepo() {
	repo := repository.NewURLRepository(db)
	fmt.Println(repo)
}

func initService() {
	service := service.NewURLService(*repo, log)
	fmt.Println(service)
}

func initHandler() {
	urlHandler := handler.NewURLHandler(db, log)
	fmt.Println(urlHandler)
}
