package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"mini_quicko/internal/models"
	"mini_quicko/internal/storage"
	"sort"
	"time"
)

type AnalysisServiceImpl struct {
	kaspiService     KaspiService
	priceHistoryRepo storage.PriceHistoryRepo
	offersCacheRepo  storage.OffersCacheRepo
	logger           *slog.Logger
}

func NewAnalysisService(kaspiService KaspiService, priceHistoryRepo storage.PriceHistoryRepo, offersCacheRepo storage.OffersCacheRepo, logger *slog.Logger) *AnalysisServiceImpl {
	return &AnalysisServiceImpl{
		kaspiService:     kaspiService,
		priceHistoryRepo: priceHistoryRepo,
		offersCacheRepo:  offersCacheRepo,
		logger:           logger,
	}
}

func (s *AnalysisServiceImpl) AnalyzeProduct(ctx context.Context, productID, cityID string) (*models.PriceAnalysis, error) {
	s.logger.Info("starting product analysis", "product_id", productID)


	cachedOffers, err := s.fetchOffersWithCache(ctx, productID, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch offers: %w", err)
	}

	if len(cachedOffers) == 0 {
		s.logger.Warn("no offers found", "product_id", productID)
		return nil, fmt.Errorf("no offers found for product %s", productID)
	}

	s.logger.Info("analyzing prices", "product_id", productID, "offers_count", len(cachedOffers))

	analysis := &models.PriceAnalysis{
		ProductID:   productID,
		TotalOffers: len(cachedOffers),
		Timestamp:   time.Now(),
	}


	prices := make([]float64, len(cachedOffers))
	var sum float64
	analysis.MinPrice = math.MaxFloat64
	analysis.MaxPrice = 0

	for i, offer := range cachedOffers {
		price := offer.Price
		prices[i] = price
		sum += price

		if price < analysis.MinPrice {
			analysis.MinPrice = price
		}
		if price > analysis.MaxPrice {
			analysis.MaxPrice = price
		}
	}

	analysis.AvgPrice = sum / float64(len(prices))


	sort.Float64s(prices)
	mid := len(prices) / 2
	if len(prices)%2 == 0 {
		analysis.MedianPrice = (prices[mid-1] + prices[mid]) / 2
	} else {
		analysis.MedianPrice = prices[mid]
	}

	// просто вкид, мб другая логика нужна
	analysis.OptimalPrice = analysis.MedianPrice * 1.05
	if analysis.OptimalPrice > analysis.MaxPrice {
		analysis.OptimalPrice = analysis.MedianPrice
	}

	// дефолт дампинг threshold - 10% ниже средней цены
	dumpingThreshold := analysis.AvgPrice * 0.9
	for _, offer := range cachedOffers {
		if offer.Price < dumpingThreshold {
			belowPercent := ((analysis.AvgPrice - offer.Price) / analysis.AvgPrice) * 100
			analysis.DumpingMerchants = append(analysis.DumpingMerchants, models.DumpingMerchant{
				MerchantID:      offer.MerchantID,
				MerchantName:    offer.MerchantName,
				Price:           offer.Price,
				BelowAvgPercent: belowPercent,
				Rating:          offer.MerchantRating,
				ReviewsCount:    offer.MerchantReviewsCount,
			})
		}
	}

	productName := ""
	if len(cachedOffers) > 0 {
		productName = cachedOffers[0].Title
	}

	history := &models.PriceHistory{
		ProductID:   productID,
		ProductName: productName,
		Timestamp:   time.Now(),
		MinPrice:    analysis.MinPrice,
		MaxPrice:    analysis.MaxPrice,
		AvgPrice:    analysis.AvgPrice,
		OfferCount:  len(cachedOffers),
	}

	if err := s.priceHistoryRepo.SavePriceHistory(ctx, history); err != nil {
		s.logger.Error("failed to save price history", "error", err, "product_id", productID)
	} else {
		s.logger.Info("price history saved", "product_id", productID)
	}

	s.logger.Info("analysis completed", "product_id", productID, "dumping_merchants", len(analysis.DumpingMerchants))
	return analysis, nil
}

func (s *AnalysisServiceImpl) GetProductSellers(ctx context.Context, productID, cityID string) ([]models.Seller, error) {
	s.logger.Info("fetching product sellers", "product_id", productID, "city_id", cityID)


	cachedOffers, err := s.fetchOffersWithCache(ctx, productID, cityID)
	if err != nil {
		s.logger.Error("failed to fetch offers", "error", err, "product_id", productID)
		return nil, err
	}

	sellers := make([]models.Seller, 0, len(cachedOffers))
	for _, offer := range cachedOffers {
		seller := models.Seller{
			MerchantID:   offer.MerchantID,
			MerchantName: offer.MerchantName,
			Rating:       offer.MerchantRating,
			ReviewsCount: offer.MerchantReviewsCount,
			Price:        offer.Price,
		}
		sellers = append(sellers, seller)
	}

	s.logger.Info("sellers fetched successfully", "product_id", productID, "count", len(sellers))
	return sellers, nil
}

func (s *AnalysisServiceImpl) GetDumpingSellers(ctx context.Context, productID, cityID string) ([]models.Seller, float64, error) {
	s.logger.Info("fetching dumping sellers", "product_id", productID, "city_id", cityID)


	cachedOffers, err := s.fetchOffersWithCache(ctx, productID, cityID)
	if err != nil {
		s.logger.Error("failed to fetch offers", "error", err, "product_id", productID)
		return nil, 0, err
	}

	if len(cachedOffers) == 0 {
		s.logger.Warn("no offers found", "product_id", productID)
		return nil, 0, fmt.Errorf("no offers found for product %s", productID)
	}


	var sum float64
	for _, offer := range cachedOffers {
		sum += offer.Price
	}
	avgPrice := sum / float64(len(cachedOffers))


	dumpingThreshold := avgPrice * 0.9
	dumpingSellers := make([]models.Seller, 0)

	for _, offer := range cachedOffers {
		if offer.Price < dumpingThreshold {
			seller := models.Seller{
				MerchantID:   offer.MerchantID,
				MerchantName: offer.MerchantName,
				Rating:       offer.MerchantRating,
				ReviewsCount: offer.MerchantReviewsCount,
				Price:        offer.Price,
			}
			dumpingSellers = append(dumpingSellers, seller)
		}
	}

	s.logger.Info("dumping sellers fetched", "product_id", productID, "count", len(dumpingSellers), "avg_price", avgPrice)
	return dumpingSellers, avgPrice, nil
}
