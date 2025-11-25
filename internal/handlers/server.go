package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"mini_quicko/internal/config"

	"github.com/gin-gonic/gin"
)


type Server struct {
	config     *config.Config
	router     *gin.Engine
	handler    *HandlerImpl
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(cfg *config.Config, handler *HandlerImpl, logger *slog.Logger) *Server {
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	handler.RegisterRoutes(router)

	return &Server{
		config:  cfg,
		router:  router,
		handler: handler,
		logger:  logger,
	}
}

func (s *Server) Start() error {
	addr := s.config.Server.GetAddress()

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.logger.Info("starting server", "address", addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	s.logger.Info("server stopped")

	return nil
}
