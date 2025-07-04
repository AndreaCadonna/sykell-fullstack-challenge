package dto

import (
	"fmt"
	"net/url"
	"strings"
)

// AddURLRequest represents a request to add a new URL for crawling
type AddURLRequest struct {
	URL string `json:"url" binding:"required"`
}

// Validate validates the URL format
func (r *AddURLRequest) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return fmt.Errorf("url cannot be empty")
	}
	
	// Parse URL to validate format
	parsedURL, err := url.Parse(r.URL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}
	
	// Ensure scheme is present (http or https)
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include protocol (http:// or https://)")
	}
	
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https protocol")
	}
	
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a valid host")
	}
	
	return nil
}

// Normalize normalizes the URL (lowercase scheme and host, remove trailing slash)
func (r *AddURLRequest) Normalize() {
	r.URL = strings.TrimSpace(r.URL)
	
	// Parse and rebuild URL for normalization
	if parsedURL, err := url.Parse(r.URL); err == nil {
		parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
		parsedURL.Host = strings.ToLower(parsedURL.Host)
		
		// Remove trailing slash unless it's the root path
		if parsedURL.Path != "/" && strings.HasSuffix(parsedURL.Path, "/") {
			parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
		}
		
		r.URL = parsedURL.String()
	}
}

// ValidateTokenRequest represents a request to validate an API token
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// PaginationRequest represents common pagination parameters
type PaginationRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
	Search   string `form:"search"`
	Status   string `form:"status"`
	SortBy   string `form:"sort_by,default=created_at"`
	SortDir  string `form:"sort_dir,default=desc"`
}

// Validate validates pagination parameters
func (p *PaginationRequest) Validate() error {
	// Validate sort direction
	if p.SortDir != "asc" && p.SortDir != "desc" {
		return fmt.Errorf("sort_dir must be 'asc' or 'desc'")
	}
	
	// Validate sort field
	validSortFields := []string{"id", "url", "status", "created_at", "updated_at"}
	isValidSort := false
	for _, field := range validSortFields {
		if p.SortBy == field {
			isValidSort = true
			break
		}
	}
	if !isValidSort {
		return fmt.Errorf("sort_by must be one of: %s", strings.Join(validSortFields, ", "))
	}
	
	// Validate status filter
	if p.Status != "" {
		validStatuses := []string{"queued", "running", "completed", "error"}
		isValidStatus := false
		for _, status := range validStatuses {
			if p.Status == status {
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			return fmt.Errorf("status must be one of: %s", strings.Join(validStatuses, ", "))
		}
	}
	
	return nil
}

// GetOffset calculates the database offset for pagination
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetOrderClause returns the ORDER BY clause for the database query
func (p *PaginationRequest) GetOrderClause() string {
	return fmt.Sprintf("%s %s", p.SortBy, strings.ToUpper(p.SortDir))
}