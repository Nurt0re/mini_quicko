package models

import "time"

type DeliveryOptions struct {
	Delivery          string  `json:"delivery"`
	KdPickupDate      string  `json:"kdPickupDate"`
	DeliveryType      string  `json:"deliveryType"`
	LocatedInPoint    string  `json:"locatedInPoint"`
	KaspiDelivery     bool    `json:"kaspiDelivery"`
	InterCity         bool    `json:"interCity"`
	DeliveryCost      float64 `json:"deliveryCost,omitempty"`
	DeliveryThreshold float64 `json:"deliveryThreshold,omitempty"`
}

type Offer struct {
	MasterSku               string                     `json:"masterSku"`
	MasterCategory          string                     `json:"masterCategory"`
	MerchantID              string                     `json:"merchantId"`
	MerchantName            string                     `json:"merchantName"`
	MerchantSku             string                     `json:"merchantSku"`
	MerchantReviewsQuantity int                        `json:"merchantReviewsQuantity"`
	MerchantRating          float64                    `json:"merchantRating"`
	Title                   string                     `json:"title"`
	Price                   float64                    `json:"price"`
	Delivery                string                     `json:"delivery"`
	KdPickupDate            string                     `json:"kdPickupDate"`
	AvailabilityDate        string                     `json:"availabilityDate"`
	KaspiDelivery           bool                       `json:"kaspiDelivery"`
	DeliveryType            string                     `json:"deliveryType"`
	DeliveryDuration        string                     `json:"deliveryDuration"`
	DeliveryOptions         map[string]DeliveryOptions `json:"deliveryOptions"`
}

type Badge struct {
	Code         string   `json:"code"`
	ImageUrl     string   `json:"imageUrl"`
	Type         string   `json:"type"`
	Priority     int      `json:"priority"`
	Owner        string   `json:"owner"`
	BonusType    string   `json:"bonusType"`
	PaymentModes []string `json:"paymentModes,omitempty"`
	Duration     int      `json:"duration,omitempty"`
}

type LoanPicker struct {
	Present                 bool `json:"present"`
	MaxInterestFreeDuration int  `json:"maxInterestFreeDuration"`
}

type ProductCardInfo struct {
	Badges []Badge `json:"badges"`
}

type KaspiResponse struct {
	Offers                   []Offer         `json:"offers"`
	Total                    int             `json:"total"`
	OffersCount              int             `json:"offersCount"`
	Badges                   []Badge         `json:"badges"`
	LoanPicker               LoanPicker      `json:"loanPicker"`
	HighRatingPresent        bool            `json:"highRatingPresent"`
	ExcellentMerchantPresent bool            `json:"excellentMerchantPresent"`
	ProductCardInfo          ProductCardInfo `json:"productCardInfo"`
}

type PriceHistory struct {
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Timestamp   time.Time `json:"timestamp"`
	MinPrice    float64   `json:"min_price"`
	MaxPrice    float64   `json:"max_price"`
	AvgPrice    float64   `json:"avg_price"`
	OfferCount  int       `json:"offer_count"`
}

type PriceAnalysis struct {
	ProductID        string            `json:"product_id"`
	MinPrice         float64           `json:"min_price"`
	MaxPrice         float64           `json:"max_price"`
	AvgPrice         float64           `json:"avg_price"`
	MedianPrice      float64           `json:"median_price"`
	OptimalPrice     float64           `json:"optimal_price"`
	TotalOffers      int               `json:"total_offers"`
	DumpingMerchants []DumpingMerchant `json:"dumping_merchants"`
	Timestamp        time.Time         `json:"timestamp"`
}

type DumpingMerchant struct {
	MerchantID      string  `json:"merchant_id"`
	MerchantName    string  `json:"merchant_name"`
	Price           float64 `json:"price"`
	BelowAvgPercent float64 `json:"below_avg_percent"`
	Rating          float64 `json:"rating"`
	ReviewsCount    int     `json:"reviews_count"`
}

type Seller struct {
	MerchantID   string  `json:"merchant_id"`
	MerchantName string  `json:"merchant_name"`
	Rating       float64 `json:"rating"`
	ReviewsCount int     `json:"reviews_count"`
	Price        float64 `json:"price"`
}

type CachedOffer struct {
	ProductID            string    `json:"product_id"`
	CityID               string    `json:"city_id"`
	MerchantID           string    `json:"merchant_id"`
	MerchantName         string    `json:"merchant_name"`
	MerchantRating       float64   `json:"merchant_rating"`
	MerchantReviewsCount int       `json:"merchant_reviews_count"`
	Title                string    `json:"title"`
	Price                float64   `json:"price"`
	TotalOffers          int       `json:"total_offers"`
	OffersCount          int       `json:"offers_count"`
	FetchedAt            time.Time `json:"fetched_at"`
}
