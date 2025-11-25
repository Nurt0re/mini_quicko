package service

import (
	"context"
	"mini_quicko/internal/models"
)

const (
	cacheMaxAgeMinutes = 30 
)


func (s *AnalysisServiceImpl) fetchOffersWithCache(ctx context.Context, productID, cityID string) ([]models.CachedOffer, error) {
	cachedOffers, err := s.offersCacheRepo.GetCachedOffers(ctx, productID, cityID, cacheMaxAgeMinutes)
	if err == nil && len(cachedOffers) > 0 {
		s.logger.Info("using cached offers", "product_id", productID, "count", len(cachedOffers))
		return cachedOffers, nil
	}


	s.logger.Info("cache miss, fetching from API", "product_id", productID)
	kaspiResp, err := s.kaspiService.FetchOffers(ctx, productID, cityID)
	if err != nil {
		return nil, err
	}


	if err := s.offersCacheRepo.SaveOffers(ctx, productID, cityID, kaspiResp.Offers, kaspiResp.Total, kaspiResp.OffersCount); err != nil {
		s.logger.Warn("failed to save offers to cache", "error", err)
	}


	cachedOffers = make([]models.CachedOffer, len(kaspiResp.Offers))
	for i, offer := range kaspiResp.Offers {
		cachedOffers[i] = models.CachedOffer{
			ProductID:            productID,
			CityID:               cityID,
			MerchantID:           offer.MerchantID,
			MerchantName:         offer.MerchantName,
			MerchantRating:       offer.MerchantRating,
			MerchantReviewsCount: offer.MerchantReviewsQuantity,
			Title:                offer.Title,
			Price:                offer.Price,
			TotalOffers:          kaspiResp.Total,
			OffersCount:          kaspiResp.OffersCount,
		}
	}

	return cachedOffers, nil
}
