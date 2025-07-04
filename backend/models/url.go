package models

import (
	"time"
	"gorm.io/gorm"
)

// URLStatus represents the status of a URL crawling process
type URLStatus string

const (
	StatusQueued    URLStatus = "queued"
	StatusRunning   URLStatus = "running"
	StatusCompleted URLStatus = "completed"
	StatusError     URLStatus = "error"
)

// URL represents a target URL for crawling
type URL struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	URL          string    `json:"url" gorm:"type:varchar(2048);not null;uniqueIndex:unique_url,length:255"`
	Status       URLStatus `json:"status" gorm:"type:enum('queued','running','completed','error');default:'queued';index"`
	ErrorMessage *string   `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"index"`
	UpdatedAt    time.Time `json:"updated_at"`
	
	// Relationships
	CrawlResult *CrawlResult `json:"crawl_result,omitempty" gorm:"foreignKey:URLID"`
	FoundLinks  []FoundLink  `json:"found_links,omitempty" gorm:"foreignKey:URLID"`
}

// TableName overrides the table name
func (URL) TableName() string {
	return "urls"
}

// BeforeCreate hook to set default status
func (u *URL) BeforeCreate(tx *gorm.DB) error {
	if u.Status == "" {
		u.Status = StatusQueued
	}
	return nil
}