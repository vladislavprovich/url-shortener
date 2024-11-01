package main

import (
	"context"
	"database/sql"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/vladislavprovich/url-shortener/internal/middleware"
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
	"github.com/vladislavprovich/url-shortener/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Initialize
	logger := initLogger()
	db := initDB()
	defer func() {

		if err := db.Close(); err != nil {

		}
	}()
	repo := initRepo(db)
	service := initService(&repo, logger)
	urlHandler := initHandler(service, logger)

	logger.Debug("main, db initialized", zap.Any("db", db), zap.Any("repo", repo))
	logger.Debug("main, logger initialized", zap.Any("logger", logger))
	logger.Debug("main, repo initialized", zap.Any("repo", repo))
	logger.Debug("main, service initialized", zap.Any("service", service))
	logger.Debug("main, URL handler initialized", zap.Any("URL_handler", urlHandler))

	// Create Router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	// Apply Middlewares
	r.Use(middleware.Recoverer(logger)) // is working
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimiter)

	// Routes
	r.Post("/shorten", urlHandler.ShortenURL)       // is working
	r.Get("/{shortURL}", urlHandler.Redirect)       // is working
	r.Get("/{shortURL}/stats", urlHandler.GetStats) // is working

	// Start Server Graceful Shutdown
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
		logger.Debug("setting server port", zap.String("port", serverPort))
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

func initLogger() *zap.Logger {
	return logger.NewLogger(os.Getenv("LOG_LEVEL"))
}

func initDB() *sql.DB {
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

func initHandler(serv service.URLService, logger *zap.Logger) *handler.URLHandler {
	return handler.NewURLHandler(serv, logger)
}
