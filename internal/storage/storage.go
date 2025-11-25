package storage

import (
	"context"
	"database/sql"
	"log/slog"
	"mini_quicko/internal/models"
)

type PriceHistoryRepo interface {
	SavePriceHistory(ctx context.Context, history *models.PriceHistory) error
	GetPriceHistory(ctx context.Context, productID string) ([]*models.PriceHistory, error)
	Close() error
}

type OffersCacheRepo interface {
	SaveOffers(ctx context.Context, productID, cityID string, offers []models.Offer, totalOffers, offersCount int) error
	GetCachedOffers(ctx context.Context, productID, cityID string, maxAge int) ([]models.CachedOffer, error)
	ClearCache(ctx context.Context, productID, cityID string) error
}

type Repository struct {
	PriceHistoryRepo PriceHistoryRepo
	OffersCacheRepo  OffersCacheRepo
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	postgresStorage := NewPostgresStorage(db, logger)
	return &Repository{
		PriceHistoryRepo: postgresStorage,
		OffersCacheRepo:  postgresStorage,
	}
}
