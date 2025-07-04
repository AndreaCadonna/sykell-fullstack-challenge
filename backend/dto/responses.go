package dto

import (
	"time"
	"web-crawler/models"
)

// APIResponse represents a standard API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta represents response metadata (pagination, etc.)
type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SuccessResponse creates a successful API response
func SuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

// PaginatedResponse creates a successful paginated API response
func PaginatedResponse(data interface{}, page, pageSize, total int) APIResponse {
	totalPages := (total + pageSize - 1) / pageSize // Ceiling division
	
	return APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(code, message, details string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// URLResponse represents a URL in API responses
type URLResponse struct {
	ID           uint                 `json:"id"`
	URL          string               `json:"url"`
	Status       models.URLStatus     `json:"status"`
	ErrorMessage *string              `json:"error_message,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	CrawlResult  *CrawlResultResponse `json:"crawl_result,omitempty"`
}

// CrawlResultResponse represents crawl results in API responses
type CrawlResultResponse struct {
	ID                     uint                   `json:"id"`
	HTMLVersion            *string                `json:"html_version"`
	PageTitle              *string                `json:"page_title"`
	HeadingCounts          map[string]int         `json:"heading_counts"`
	InternalLinksCount     int                    `json:"internal_links_count"`
	ExternalLinksCount     int                    `json:"external_links_count"`
	InaccessibleLinksCount int                    `json:"inaccessible_links_count"`
	HasLoginForm           bool                   `json:"has_login_form"`
	CrawledAt              time.Time              `json:"crawled_at"`
	CrawlDurationMs        *int                   `json:"crawl_duration_ms"`
	TotalLinks             int                    `json:"total_links"`
}

// FoundLinkResponse represents a found link in API responses
type FoundLinkResponse struct {
	ID           uint    `json:"id"`
	LinkURL      string  `json:"link_url"`
	LinkText     *string `json:"link_text"`
	IsInternal   bool    `json:"is_internal"`
	IsAccessible *bool   `json:"is_accessible"`
	StatusCode   *int    `json:"status_code"`
	ErrorMessage *string `json:"error_message"`
	IsBroken     bool    `json:"is_broken"`
	StatusCategory string `json:"status_category"`
	CreatedAt    time.Time `json:"created_at"`
}

// URLDetailResponse represents detailed URL information with links
type URLDetailResponse struct {
	URLResponse
	FoundLinks []FoundLinkResponse `json:"found_links"`
}

// TokenValidationResponse represents token validation response
type TokenValidationResponse struct {
	Valid     bool       `json:"valid"`
	TokenName string     `json:"token_name,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// FromURL converts a models.URL to URLResponse
func FromURL(url *models.URL) URLResponse {
	response := URLResponse{
		ID:           url.ID,
		URL:          url.URL,
		Status:       url.Status,
		ErrorMessage: url.ErrorMessage,
		CreatedAt:    url.CreatedAt,
		UpdatedAt:    url.UpdatedAt,
	}
	
	// Include crawl result if available
	if url.CrawlResult != nil {
		response.CrawlResult = FromCrawlResult(url.CrawlResult)
	}
	
	return response
}

// FromCrawlResult converts a models.CrawlResult to CrawlResultResponse
func FromCrawlResult(result *models.CrawlResult) *CrawlResultResponse {
	return &CrawlResultResponse{
		ID:                     result.ID,
		HTMLVersion:            result.HTMLVersion,
		PageTitle:              result.PageTitle,
		HeadingCounts:          result.GetHeadingCounts(),
		InternalLinksCount:     result.InternalLinksCount,
		ExternalLinksCount:     result.ExternalLinksCount,
		InaccessibleLinksCount: result.InaccessibleLinksCount,
		HasLoginForm:           result.HasLoginForm,
		CrawledAt:              result.CrawledAt,
		CrawlDurationMs:        result.CrawlDurationMs,
		TotalLinks:             result.GetTotalLinks(),
	}
}

// FromFoundLink converts a models.FoundLink to FoundLinkResponse
func FromFoundLink(link *models.FoundLink) FoundLinkResponse {
	return FoundLinkResponse{
		ID:             link.ID,
		LinkURL:        link.LinkURL,
		LinkText:       link.LinkText,
		IsInternal:     link.IsInternal,
		IsAccessible:   link.IsAccessible,
		StatusCode:     link.StatusCode,
		ErrorMessage:   link.ErrorMessage,
		IsBroken:       link.IsBroken(),
		StatusCategory: link.GetStatusCategory(),
		CreatedAt:      link.CreatedAt,
	}
}

// FromURLs converts a slice of models.URL to slice of URLResponse
func FromURLs(urls []models.URL) []URLResponse {
	responses := make([]URLResponse, len(urls))
	for i, url := range urls {
		responses[i] = FromURL(&url)
	}
	return responses
}

// FromFoundLinks converts a slice of models.FoundLink to slice of FoundLinkResponse
func FromFoundLinks(links []models.FoundLink) []FoundLinkResponse {
	responses := make([]FoundLinkResponse, len(links))
	for i, link := range links {
		responses[i] = FromFoundLink(&link)
	}
	return responses
}