package handlers

import (
	"log/slog"
	"mini_quicko/internal/config"
	"mini_quicko/internal/service"
)

type Handler interface {
	RegisterRoutes()
}

type HandlerImpl struct {
	config  *config.Config
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(cfg *config.Config, svc *service.Service, logger *slog.Logger) *HandlerImpl {
	return &HandlerImpl{
		config:  cfg,
		service: svc,
		logger:  logger,
	}
}
