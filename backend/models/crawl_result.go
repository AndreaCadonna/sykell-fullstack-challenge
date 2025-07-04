package models

import (
	"time"

	"gorm.io/gorm"
)

// CrawlResult stores the extracted data from crawling a URL
type CrawlResult struct {
	ID    uint `json:"id" gorm:"primaryKey"`
	URLID uint `json:"url_id" gorm:"not null;index"`

	// Extracted crawl data
	HTMLVersion            *string `json:"html_version" gorm:"type:varchar(50)"`
	PageTitle              *string `json:"page_title" gorm:"type:varchar(500)"`
	H1Count                int     `json:"h1_count" gorm:"default:0"`
	H2Count                int     `json:"h2_count" gorm:"default:0"`
	H3Count                int     `json:"h3_count" gorm:"default:0"`
	H4Count                int     `json:"h4_count" gorm:"default:0"`
	H5Count                int     `json:"h5_count" gorm:"default:0"`
	H6Count                int     `json:"h6_count" gorm:"default:0"`
	InternalLinksCount     int     `json:"internal_links_count" gorm:"default:0"`
	ExternalLinksCount     int     `json:"external_links_count" gorm:"default:0"`
	InaccessibleLinksCount int     `json:"inaccessible_links_count" gorm:"default:0"`
	HasLoginForm           bool    `json:"has_login_form" gorm:"default:false"`

	// Metadata
	CrawledAt       time.Time `json:"crawled_at" gorm:"index"`
	CrawlDurationMs *int      `json:"crawl_duration_ms"`

	// Relationships
	URL *URL `json:"url,omitempty" gorm:"foreignKey:URLID"`
}

// TableName overrides the table name
func (CrawlResult) TableName() string {
	return "crawl_results"
}

// GetHeadingCounts returns a map of heading counts
func (cr *CrawlResult) GetHeadingCounts() map[string]int {
	return map[string]int{
		"h1": cr.H1Count,
		"h2": cr.H2Count,
		"h3": cr.H3Count,
		"h4": cr.H4Count,
		"h5": cr.H5Count,
		"h6": cr.H6Count,
	}
}

// GetTotalLinks returns total internal + external links
func (cr *CrawlResult) GetTotalLinks() int {
	return cr.InternalLinksCount + cr.ExternalLinksCount
}

// BeforeCreate hook to set crawled_at timestamp
func (cr *CrawlResult) BeforeCreate(tx *gorm.DB) error {
	if cr.CrawledAt.IsZero() {
		cr.CrawledAt = time.Now().UTC()
	}
	return nil
}
