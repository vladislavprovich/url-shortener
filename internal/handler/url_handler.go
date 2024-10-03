package handler

import (
	"database/sql"
	"github.com/vladislavprovich/url-shortener/internal/repository"
	"github.com/vladislavprovich/url-shortener/internal/service"
	"go.uber.org/zap"
)

type URLHandler struct {
	service service.URLService
	logger  *zap.Logger
}

func NewURLHandler(db *sql.DB, logger *zap.Logger) *URLHandler {
	repo := repository.NewURLRepository(db)
	urlService := service.NewURLService(repo)

	return &URLHandler{
		service: urlService,
		logger:  logger,
	}
}
