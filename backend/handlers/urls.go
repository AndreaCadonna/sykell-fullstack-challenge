package handlers

import (
	"net/http"
	"strconv"

	"web-crawler/database"
	"web-crawler/dto"
	"web-crawler/models"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// URLHandler handles URL-related requests
type URLHandler struct{}

// NewURLHandler creates a new URL handler
func NewURLHandler() *URLHandler {
	return &URLHandler{}
}

// ListURLs returns a paginated list of URLs
// GET /api/urls
func (h *URLHandler) ListURLs(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_PARAMS",
			"Invalid query parameters",
			err.Error(),
		))
		return
	}
	
	// Validate pagination parameters
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_PARAMS",
			"Invalid pagination parameters",
			err.Error(),
		))
		return
	}
	
	// Build query
	query := database.DB.Model(&models.URL{}).Preload("CrawlResult")
	
	// Apply filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		query = query.Where("url LIKE ?", searchPattern)
	}
	
	// Get total count for pagination
	var total int64
	if err := query.Count(&total); err.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to count URLs",
			err.Error.Error(),
		))
		return
	}
	
	// Apply pagination and sorting
	var urls []models.URL
	result := query.
		Order(req.GetOrderClause()).
		Offset(req.GetOffset()).
		Limit(req.PageSize).
		Find(&urls)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch URLs",
			result.Error.Error(),
		))
		return
	}
	
	// Convert to response format
	urlResponses := dto.FromURLs(urls)
	
	c.JSON(http.StatusOK, dto.PaginatedResponse(
		urlResponses,
		req.Page,
		req.PageSize,
		int(total),
	))
}

// AddURL adds a new URL for crawling
// POST /api/urls
func (h *URLHandler) AddURL(c *gin.Context) {
	var req dto.AddURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}
	
	// Validate URL format
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_URL",
			"Invalid URL format",
			err.Error(),
		))
		return
	}
	
	// Normalize URL
	req.Normalize()
	
	// Check if URL already exists
	var existingURL models.URL
	result := database.DB.Where("url = ?", req.URL).First(&existingURL)
	if result.Error == nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse(
			"URL_EXISTS",
			"URL already exists in the system",
			"",
		))
		return
	}
	
	// Create new URL
	newURL := models.URL{
		URL:    req.URL,
		Status: models.StatusQueued,
	}
	
	if err := database.DB.Create(&newURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to save URL",
			err.Error(),
		))
		return
	}
	
	// Return created URL
	c.JSON(http.StatusCreated, dto.SuccessResponse(dto.FromURL(&newURL)))
}

// GetURL returns details of a specific URL
// GET /api/urls/:id
func (h *URLHandler) GetURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid URL ID",
			"ID must be a positive integer",
		))
		return
	}
	
	var url models.URL
	result := database.DB.Preload("CrawlResult").First(&url, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"URL_NOT_FOUND",
				"URL not found",
				"",
			))
			return
		}
		
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch URL",
			result.Error.Error(),
		))
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse(dto.FromURL(&url)))
}

// GetURLDetails returns detailed URL information including found links
// GET /api/urls/:id/details
func (h *URLHandler) GetURLDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid URL ID",
			"ID must be a positive integer",
		))
		return
	}
	
	var url models.URL
	result := database.DB.
		Preload("CrawlResult").
		Preload("FoundLinks").
		First(&url, id)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"URL_NOT_FOUND",
				"URL not found",
				"",
			))
			return
		}
		
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch URL details",
			result.Error.Error(),
		))
		return
	}
	
	// Build detailed response
	response := dto.URLDetailResponse{
		URLResponse: dto.FromURL(&url),
		FoundLinks:  dto.FromFoundLinks(url.FoundLinks),
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

// DeleteURL deletes a URL and all related data
// DELETE /api/urls/:id
func (h *URLHandler) DeleteURL(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid URL ID",
			"ID must be a positive integer",
		))
		return
	}
	
	// Check if URL exists
	var url models.URL
	result := database.DB.First(&url, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				"URL_NOT_FOUND",
				"URL not found",
				"",
			))
			return
		}
		
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to check URL existence",
			result.Error.Error(),
		))
		return
	}
	
	// Delete URL (cascading delete will handle related records)
	if err := database.DB.Delete(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to delete URL",
			err.Error(),
		))
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"message": "URL deleted successfully",
		"id":      id,
	}))
}

// BulkDeleteURLs deletes multiple URLs
// DELETE /api/urls/bulk
func (h *URLHandler) BulkDeleteURLs(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}
	
	// Delete URLs in transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	result := tx.Where("id IN ?", req.IDs).Delete(&models.URL{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to delete URLs",
			result.Error.Error(),
		))
		return
	}
	
	tx.Commit()
	
	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"message":       "URLs deleted successfully",
		"deleted_count": result.RowsAffected,
		"ids":           req.IDs,
	}))
}