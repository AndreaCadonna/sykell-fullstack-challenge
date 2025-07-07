package handlers

import (
	"net/http"
	"strconv"

	"web-crawler/database"
	"web-crawler/dto"
	"web-crawler/models"
	"web-crawler/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CrawlHandler handles crawl-related requests
type CrawlHandler struct {
	crawlManager *services.CrawlManager
}

// NewCrawlHandler creates a new crawl handler
func NewCrawlHandler(crawlManager *services.CrawlManager) *CrawlHandler {
	return &CrawlHandler{
		crawlManager: crawlManager,
	}
}

// StartCrawl starts crawling a specific URL
// POST /api/urls/:id/crawl
func (h *CrawlHandler) StartCrawl(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid URL ID",
			"ID must be a positive integer",
		))
		return
	}

	// Get URL from database
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
			"Failed to fetch URL",
			result.Error.Error(),
		))
		return
	}

	// Check if URL is in a valid state for crawling
	if url.Status == models.StatusRunning {
		c.JSON(http.StatusConflict, dto.ErrorResponse(
			"CRAWL_IN_PROGRESS",
			"URL is already being crawled",
			"Wait for current crawl to complete before starting a new one",
		))
		return
	}

	// Queue the URL for crawling
	if err := h.crawlManager.QueueURL(uint(id), url.URL); err != nil {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse(
			"QUEUE_FULL",
			"Crawl queue is full",
			err.Error(),
		))
		return
	}

	// Update URL status to queued (if not already)
	if url.Status != models.StatusQueued {
		updates := map[string]interface{}{
			"status":        models.StatusQueued,
			"error_message": nil,
		}
		database.DB.Model(&url).Updates(updates)
	}

	c.JSON(http.StatusAccepted, dto.SuccessResponse(gin.H{
		"message":    "Crawl started successfully",
		"url_id":     id,
		"url":        url.URL,
		"status":     "queued",
		"queue_info": h.crawlManager.GetQueueStatus(),
	}))
}

// GetCrawlStatus returns the current status of a URL crawl
// GET /api/urls/:id/crawl/status
func (h *CrawlHandler) GetCrawlStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_ID",
			"Invalid URL ID",
			"ID must be a positive integer",
		))
		return
	}

	// Get URL with crawl result
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

	// Build status response
	statusResponse := gin.H{
		"url_id":     id,
		"url":        url.URL,
		"status":     url.Status,
		"created_at": url.CreatedAt,
		"updated_at": url.UpdatedAt,
		"queue_info": h.crawlManager.GetQueueStatus(),
	}

	if url.ErrorMessage != nil {
		statusResponse["error_message"] = *url.ErrorMessage
	}

	// Include crawl result if available
	if url.CrawlResult != nil {
		statusResponse["crawl_result"] = dto.FromCrawlResult(url.CrawlResult)
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(statusResponse))
}

// StartBulkCrawl starts crawling multiple URLs
// POST /api/crawls/bulk
func (h *CrawlHandler) StartBulkCrawl(c *gin.Context) {
	var req struct {
		URLIDs []uint `json:"url_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}

	// Limit bulk operation size
	const maxBulkSize = 10
	if len(req.URLIDs) > maxBulkSize {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"TOO_MANY_URLS",
			"Too many URLs in bulk request",
			"Maximum 10 URLs allowed per bulk operation",
		))
		return
	}

	// Get URLs from database
	var urls []models.URL
	result := database.DB.Where("id IN ?", req.URLIDs).Find(&urls)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"DATABASE_ERROR",
			"Failed to fetch URLs",
			result.Error.Error(),
		))
		return
	}

	// Check if we found all requested URLs
	if len(urls) != len(req.URLIDs) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"URLS_NOT_FOUND",
			"Some URLs not found",
			"One or more URL IDs do not exist",
		))
		return
	}

	// Queue URLs for crawling
	queueResults := make([]gin.H, 0, len(urls))
	successCount := 0

	for _, url := range urls {
		result := gin.H{
			"url_id": url.ID,
			"url":    url.URL,
		}

		// Check if URL is already being crawled
		if url.Status == models.StatusRunning {
			result["status"] = "skipped"
			result["reason"] = "already in progress"
		} else {
			// Try to queue the URL
			if err := h.crawlManager.QueueURL(url.ID, url.URL); err != nil {
				result["status"] = "failed"
				result["reason"] = err.Error()
			} else {
				// Update URL status to queued
				if url.Status != models.StatusQueued {
					updates := map[string]interface{}{
						"status":        models.StatusQueued,
						"error_message": nil,
					}
					database.DB.Model(&url).Updates(updates)
				}

				result["status"] = "queued"
				successCount++
			}
		}

		queueResults = append(queueResults, result)
	}

	c.JSON(http.StatusAccepted, dto.SuccessResponse(gin.H{
		"message":       "Bulk crawl operation completed",
		"total_urls":    len(urls),
		"queued_count":  successCount,
		"skipped_count": len(urls) - successCount,
		"results":       queueResults,
		"queue_info":    h.crawlManager.GetQueueStatus(),
	}))
}

// GetQueueStatus returns information about the crawl queue
// GET /api/crawls/queue/status
func (h *CrawlHandler) GetQueueStatus(c *gin.Context) {
	queueStatus := h.crawlManager.GetQueueStatus()

	// Add additional statistics
	var stats struct {
		QueuedCount    int64 `json:"queued_count"`
		RunningCount   int64 `json:"running_count"`
		CompletedCount int64 `json:"completed_count"`
		ErrorCount     int64 `json:"error_count"`
	}

	database.DB.Model(&models.URL{}).Where("status = ?", models.StatusQueued).Count(&stats.QueuedCount)
	database.DB.Model(&models.URL{}).Where("status = ?", models.StatusRunning).Count(&stats.RunningCount)
	database.DB.Model(&models.URL{}).Where("status = ?", models.StatusCompleted).Count(&stats.CompletedCount)
	database.DB.Model(&models.URL{}).Where("status = ?", models.StatusError).Count(&stats.ErrorCount)

	response := gin.H{
		"queue_manager":  queueStatus,
		"database_stats": stats,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
