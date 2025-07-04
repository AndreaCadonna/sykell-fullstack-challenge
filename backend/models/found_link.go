package models

import (
	"time"
)

// FoundLink represents a link discovered during crawling
type FoundLink struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	URLID        uint    `json:"url_id" gorm:"not null;index"`
	LinkURL      string  `json:"link_url" gorm:"type:varchar(2048);not null"`
	LinkText     *string `json:"link_text" gorm:"type:varchar(500)"`
	IsInternal   bool    `json:"is_internal" gorm:"not null;index"`
	IsAccessible *bool   `json:"is_accessible" gorm:"index"` // NULL = not checked yet
	StatusCode   *int    `json:"status_code" gorm:"index"`
	ErrorMessage *string `json:"error_message" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"`
	
	// Relationships
	URL *URL `json:"url,omitempty" gorm:"foreignKey:URLID"`
}

// TableName overrides the table name
func (FoundLink) TableName() string {
	return "found_links"
}

// IsBroken returns true if the link is inaccessible (4xx or 5xx status codes)
func (fl *FoundLink) IsBroken() bool {
	if fl.StatusCode == nil {
		return false
	}
	return *fl.StatusCode >= 400
}

// GetStatusCategory returns a human-readable status category
func (fl *FoundLink) GetStatusCategory() string {
	if fl.StatusCode == nil {
		return "unchecked"
	}
	
	code := *fl.StatusCode
	switch {
	case code >= 200 && code < 300:
		return "success"
	case code >= 300 && code < 400:
		return "redirect"
	case code >= 400 && code < 500:
		return "client_error"
	case code >= 500:
		return "server_error"
	default:
		return "unknown"
	}
}