package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/internal/service"
	"github.com/vladislavprovich/url-shortener/internal/validator"
	"go.uber.org/zap"
)

type URLHandler struct {
	service service.URLService
	logger  *zap.Logger
}

// TODO repository and service must init in main
func NewURLHandler(db *sql.DB, logger *zap.Logger) *URLHandler {
	repo := repository.NewURLRepository(db)
	service := service.NewURLService(repo, logger)
	return &URLHandler{
		service: service,
		logger:  logger,
	}
}

// TODO must be in storage layer
func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Verify connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handler.ShortenURL called")
	var req models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("handler, failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate Request
	if err := validator.Validate(req); err != nil {
		h.logger.Warn("handler, validation failed", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create Short URL
	shortURL, err := h.service.CreateShortURL(r.Context(), req)
	if err != nil {
		h.logger.Error("handler, failed to create short URL", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	baseURL := r.Host
	if os.Getenv("BASE_URL") != "" {
		baseURL = os.Getenv("BASE_URL")
	}

	response := models.ShortenResponse{
		ShortURL: fmt.Sprintf("%s/%s", baseURL, shortURL),
	}

	h.logger.Info("handler, short URL created", zap.String("short_url", response.ShortURL))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("handler, failed to write response", zap.Error(err))
	}
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	originalURL, err := h.service.GetOriginalURL(r.Context(), shortURL)
	if err != nil {
		// TODO add logger
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Log the redirect
	referrer := r.Referer()
	h.service.LogRedirect(r.Context(), shortURL, referrer)

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	stats, err := h.service.GetStats(r.Context(), shortURL)
	if err != nil {
		// TODO add logger
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
