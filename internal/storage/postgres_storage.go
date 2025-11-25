package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"mini_quicko/internal/models"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresStorage(db *sql.DB, logger *slog.Logger) *PostgresStorage {
	return &PostgresStorage{
		db:     db,
		logger: logger,
	}
}

func (s *PostgresStorage) SavePriceHistory(ctx context.Context, history *models.PriceHistory) error {
	s.logger.Debug("saving price history", "product_id", history.ProductID)

	query := `
		INSERT INTO price_history (product_id, product_name, timestamp, min_price, max_price, avg_price, offer_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.ExecContext(ctx, query,
		history.ProductID,
		history.ProductName,
		history.Timestamp,
		history.MinPrice,
		history.MaxPrice,
		history.AvgPrice,
		history.OfferCount,
	)

	if err != nil {
		return fmt.Errorf("failed to save price history: %w", err)
	}

	s.logger.Debug("price history saved successfully", "product_id", history.ProductID)
	return nil
}

func (s *PostgresStorage) GetPriceHistory(ctx context.Context, productID string) ([]*models.PriceHistory, error) {
	s.logger.Debug("fetching price history", "product_id", productID)

	query := `
		SELECT product_id, product_name, timestamp, min_price, max_price, avg_price, offer_count
		FROM price_history
		WHERE product_id = $1
		ORDER BY timestamp DESC
		LIMIT 100
	`

	rows, err := s.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to query price history: %w", err)
	}
	defer rows.Close()

	var histories []*models.PriceHistory
	for rows.Next() {
		var h models.PriceHistory
		if err := rows.Scan(&h.ProductID, &h.ProductName, &h.Timestamp, &h.MinPrice, &h.MaxPrice, &h.AvgPrice, &h.OfferCount); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		histories = append(histories, &h)
	}

	return histories, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func (s *PostgresStorage) SaveOffers(ctx context.Context, productID, cityID string, offers []models.Offer, totalOffers, offersCount int) error {
	s.logger.Debug("saving offers to cache", "product_id", productID, "city_id", cityID, "count", len(offers))

	if err := s.ClearCache(ctx, productID, cityID); err != nil {
		s.logger.Warn("failed to clear cache", "error", err)
	}

	query := `
		INSERT INTO offers_cache 
		(product_id, city_id, merchant_id, merchant_name, merchant_rating, merchant_reviews_count, title, price, total_offers, offers_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	for _, offer := range offers {
		_, err := s.db.ExecContext(ctx, query,
			productID,
			cityID,
			offer.MerchantID,
			offer.MerchantName,
			offer.MerchantRating,
			offer.MerchantReviewsQuantity,
			offer.Title,
			offer.Price,
			totalOffers,
			offersCount,
		)
		if err != nil {
			return fmt.Errorf("failed to save offer to cache: %w", err)
		}
	}

	s.logger.Debug("offers saved to cache successfully", "product_id", productID, "count", len(offers))
	return nil
}

func (s *PostgresStorage) GetCachedOffers(ctx context.Context, productID, cityID string, maxAgeMinutes int) ([]models.CachedOffer, error) {
	s.logger.Debug("fetching cached offers", "product_id", productID, "city_id", cityID, "max_age_minutes", maxAgeMinutes)

	query := `
		SELECT product_id, city_id, merchant_id, merchant_name, merchant_rating, merchant_reviews_count, 
		       title, price, total_offers, offers_count, fetched_at
		FROM offers_cache
		WHERE product_id = $1 AND city_id = $2 
		  AND fetched_at > NOW() - INTERVAL '1 minute' * $3
		ORDER BY price ASC
	`

	rows, err := s.db.QueryContext(ctx, query, productID, cityID, maxAgeMinutes)
	if err != nil {
		return nil, fmt.Errorf("failed to query cached offers: %w", err)
	}
	defer rows.Close()

	var cachedOffers []models.CachedOffer
	for rows.Next() {
		var co models.CachedOffer
		if err := rows.Scan(
			&co.ProductID, &co.CityID, &co.MerchantID, &co.MerchantName,
			&co.MerchantRating, &co.MerchantReviewsCount, &co.Title, &co.Price,
			&co.TotalOffers, &co.OffersCount, &co.FetchedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan cached offer: %w", err)
		}
		cachedOffers = append(cachedOffers, co)
	}

	s.logger.Debug("cached offers fetched", "product_id", productID, "count", len(cachedOffers))
	return cachedOffers, nil
}

func (s *PostgresStorage) ClearCache(ctx context.Context, productID, cityID string) error {
	s.logger.Debug("clearing cache", "product_id", productID, "city_id", cityID)

	query := `DELETE FROM offers_cache WHERE product_id = $1 AND city_id = $2`
	_, err := s.db.ExecContext(ctx, query, productID, cityID)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	return nil
}
