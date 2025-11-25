package service

import (
	"context"
	"log/slog"
	"mini_quicko/internal/models"
	"mini_quicko/internal/storage"
)

type HistoryServiceImpl struct {
	priceHistoryRepo storage.PriceHistoryRepo
	logger           *slog.Logger
}

func NewHistoryService(priceHistoryRepo storage.PriceHistoryRepo, logger *slog.Logger) *HistoryServiceImpl {
	return &HistoryServiceImpl{
		priceHistoryRepo: priceHistoryRepo,
		logger:           logger,
	}
}

func (s *HistoryServiceImpl) GetPriceHistory(ctx context.Context, productID string) ([]*models.PriceHistory, error) {
	s.logger.Info("fetching price history", "product_id", productID)

	histories, err := s.priceHistoryRepo.GetPriceHistory(ctx, productID)
	if err != nil {
		s.logger.Error("failed to get price history", "error", err, "product_id", productID)
		return nil, err
	}

	s.logger.Info("price history fetched", "product_id", productID, "count", len(histories))
	return histories, nil
}
