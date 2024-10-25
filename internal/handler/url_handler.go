package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/vladislavprovich/url-shortener/internal/models"
	"github.com/vladislavprovich/url-shortener/internal/service"
	"github.com/vladislavprovich/url-shortener/internal/validator"
	"go.uber.org/zap"
)

type URLHandler struct {
	service service.URLService
	logger  *zap.Logger
}

func NewURLHandler(srv service.URLService, logger *zap.Logger) *URLHandler {
	return &URLHandler{
		service: srv,
		logger:  logger,
	}
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
	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("handler, failed to write response", zap.Error(err))
	}
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handler.Redirect called")
	shortURL := chi.URLParam(r, "shortURL")
	originalURL, err := h.service.GetOriginalURL(r.Context(), shortURL)
	if err != nil {
		h.logger.Error("handler, failed to get original URL", zap.Error(err))
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	// Log the redirect
	referrer := r.Referer()
	h.logger.Info("handler, referrer ", zap.String("referrer", referrer))
	err = h.service.LogRedirect(r.Context(), shortURL, referrer)
	if err != nil {
		h.logger.Error("handler, failed to redirect", zap.Error(err))
		http.Error(w, "LogRedirect error", http.StatusNotFound)
		return
	}
	h.logger.Info("handler, redirect successfully")
	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handler.GetStats called")
	shortURL := chi.URLParam(r, "shortURL")
	stats, err := h.service.GetStats(r.Context(), shortURL)
	if err != nil {
		h.logger.Error("handler, failed to get stats", zap.Error(err))
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
