package service

import (
	"context"
	"log/slog"
	"mini_quicko/internal/models"
	"mini_quicko/internal/storage"
)

type KaspiService interface {
	FetchOffers(ctx context.Context, productID, cityID string) (*models.KaspiResponse, error)
}

type AnalysisService interface {
	AnalyzeProduct(ctx context.Context, productID, cityID string) (*models.PriceAnalysis, error)
	GetProductSellers(ctx context.Context, productID, cityID string) ([]models.Seller, error)
	GetDumpingSellers(ctx context.Context, productID, cityID string) ([]models.Seller, float64, error)
}

type HistoryService interface {
	GetPriceHistory(ctx context.Context, productID string) ([]*models.PriceHistory, error)
}

type Service struct {
	KaspiService    KaspiService
	AnalysisService AnalysisService
	HistoryService  HistoryService
}

func NewService(repo *storage.Repository, logger *slog.Logger) *Service {
	kaspiService := NewKaspiService()
	analysisService := NewAnalysisService(kaspiService, repo.PriceHistoryRepo, repo.OffersCacheRepo, logger)
	historyService := NewHistoryService(repo.PriceHistoryRepo, logger)

	return &Service{
		KaspiService:    kaspiService,
		AnalysisService: analysisService,
		HistoryService:  historyService,
	}
}
