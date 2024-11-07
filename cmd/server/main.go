package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/internal/repository/postgres"
	"github.com/vladislavprovich/url-shortener/internal/service"

	"github.com/vladislavprovich/url-shortener/internal/handler"
	"github.com/vladislavprovich/url-shortener/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Load Configuration
	cfg, err := LoadConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := initLogger(cfg.Logger.Level)
	db, err := postgres.PrepareConnection(ctx, cfg.Database, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Warn("Error closing db", zap.Error(err))
		}
	}()
	repo := initRepo(db)
	service := initService(&repo, logger)
	urlHandler := initHandler(service, logger, cfg.Server)
	r := handler.InitRouter(urlHandler, logger, cfg.Server)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  time.Minute,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Server is starting", zap.String("port", cfg.Server.Port))
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server error", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

func initLogger(logLevel string) *zap.Logger {
	return logger.NewLogger(logLevel)
}

func initRepo(db *sql.DB) repository.URLRepository {
	return repository.NewURLRepository(db)
}

func initService(repo *repository.URLRepository, logger *zap.Logger) service.URLService {
	return service.NewURLService(*repo, logger)
}

func initHandler(srv service.URLService, logger *zap.Logger, cfg handler.Config) *handler.URLHandler {
	return handler.NewURLHandler(srv, logger, cfg)
}
