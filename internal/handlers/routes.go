package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) RegisterRoutes(router *gin.Engine) {

	router.GET("/health", h.healthCheck)

	v1 := router.Group("/api/v1")
	{
		//анализ цен продукта
		v1.POST("/analyze", h.analyzeProduct)
		//история цен
		v1.GET("/history/:product_id", h.getHistory)
		//вся информация
		v1.GET("/products/:product_id/offers", h.getOffers)
		// информация по продавцам
		v1.GET("/products/:product_id/sellers", h.getSellers)
		// продавцы дампящие цены
		v1.GET("/products/:product_id/dumping-sellers", h.getDumpingSellers)
	}
}

func (h *HandlerImpl) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

func (h *HandlerImpl) analyzeProduct(c *gin.Context) {
	var req struct {
		ProductID string `json:"product_id" binding:"required"`
		CityID    string `json:"city_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CityID == "" {
		req.CityID = h.config.Kaspi.DefaultCityID
	}

	h.logger.Info("analyzing product", "product_id", req.ProductID, "city_id", req.CityID)

	analysis, err := h.service.AnalysisService.AnalyzeProduct(c.Request.Context(), req.ProductID, req.CityID)
	if err != nil {
		h.logger.Error("failed to analyze product", "error", err, "product_id", req.ProductID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to analyze product: %v", err)})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

func (h *HandlerImpl) getHistory(c *gin.Context) {
	productID := c.Param("product_id")

	h.logger.Info("fetching price history", "product_id", productID)

	histories, err := h.service.HistoryService.GetPriceHistory(c.Request.Context(), productID)
	if err != nil {
		h.logger.Error("failed to get price history", "error", err, "product_id", productID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get history: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"history":    histories,
	})
}

func (h *HandlerImpl) getOffers(c *gin.Context) {
	productID := c.Param("product_id")
	cityID := c.DefaultQuery("city_id", h.config.Kaspi.DefaultCityID)

	h.logger.Info("fetching offers", "product_id", productID, "city_id", cityID)

	offers, err := h.service.KaspiService.FetchOffers(c.Request.Context(), productID, cityID)
	if err != nil {
		h.logger.Error("failed to fetch offers", "error", err, "product_id", productID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch offers: %v", err)})
		return
	}

	c.JSON(http.StatusOK, offers)
}

func (h *HandlerImpl) getSellers(c *gin.Context) {
	productID := c.Param("product_id")
	cityID := c.DefaultQuery("city_id", h.config.Kaspi.DefaultCityID)

	h.logger.Info("fetching product sellers", "product_id", productID, "city_id", cityID)

	sellers, err := h.service.AnalysisService.GetProductSellers(c.Request.Context(), productID, cityID)
	if err != nil {
		h.logger.Error("failed to fetch sellers", "error", err, "product_id", productID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch sellers: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id": productID,
		"sellers":    sellers,
		"count":      len(sellers),
	})
}

func (h *HandlerImpl) getDumpingSellers(c *gin.Context) {
	productID := c.Param("product_id")
	cityID := c.DefaultQuery("city_id", h.config.Kaspi.DefaultCityID)

	h.logger.Info("fetching dumping sellers", "product_id", productID, "city_id", cityID)

	sellers, avgPrice, err := h.service.AnalysisService.GetDumpingSellers(c.Request.Context(), productID, cityID)
	if err != nil {
		h.logger.Error("failed to fetch dumping sellers", "error", err, "product_id", productID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch dumping sellers: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id":      productID,
		"average_price":   avgPrice,
		"dumping_sellers": sellers,
		"dumping_count":   len(sellers),
	})
}
