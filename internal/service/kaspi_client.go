package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"mini_quicko/internal/models"
)

const (
	defaultCityID  = "710000000"
	offersEndpoint = "https://kaspi.kz/yml/offer-view/offers/"
)

type KaspiClient struct {
	httpClient *http.Client
}

func NewKaspiService() *KaspiClient {
	return &KaspiClient{
		httpClient: &http.Client{},
	}
}


func (kc *KaspiClient) FetchOffers(ctx context.Context, productID, cityID string) (*models.KaspiResponse, error) {
	if cityID == "" {
		cityID = defaultCityID
	}
	

	url := offersEndpoint + productID

	requestBody := map[string]interface{}{
		"cityId":      cityID,
		"id":          productID,
		"merchantUID": []string{},
		"limit":       0,
		"page":        0,
		// "product": map[string]interface{}{
		// 	"brand":            "Без бренда",
		// 	"categoryCodes":    []string{"Phone cases", "Phone accessories", "Smartphones and gadgets", "Categories"},
		// 	"baseProductCodes": []string{},
		// 	"groups":           nil,
		// },
		"sortOption":          "PRICE",
		"highRating":          nil,
		"searchText":          nil,
		"isExcellentMerchant": false,
		"zoneId":              []string{"Magnum_ZONE5"},
		"installationId":      "-1",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru,en;q=0.9")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://kaspi.kz/")
	req.Header.Set("Origin", "https://kaspi.kz")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch offers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kaspi API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var kaspiResp models.KaspiResponse
	if err := json.Unmarshal(body, &kaspiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &kaspiResp, nil
}
