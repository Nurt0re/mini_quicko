package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mini_quicko/internal/config"
	"mini_quicko/internal/handlers"
	"mini_quicko/internal/service"
	"mini_quicko/internal/storage"
)

func main() {

	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		logger := setupLogger("info")
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.App.LogLevel)
	logger.Info("configuration loaded", "environment", cfg.App.Environment, "log_level", cfg.App.LogLevel)

	db, err := storage.ConnectPostgres(storage.PostgresConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connection established")

	repo := storage.NewRepository(db, logger)

	svc := service.NewService(repo, logger)

	handler := handlers.NewHandler(cfg, svc, logger)

	server := handlers.NewServer(cfg, handler, logger)

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	if err := repo.PriceHistoryRepo.Close(); err != nil {
		logger.Warn("failed to close database connection", "error", err)
	}

	logger.Info("server exited successfully")
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
