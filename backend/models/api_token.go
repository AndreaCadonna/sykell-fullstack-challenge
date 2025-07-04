package models

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// APIToken represents an authentication token
type APIToken struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	TokenHash  string     `json:"-" gorm:"type:varchar(255);not null;uniqueIndex"` // Never expose in JSON
	Name       string     `json:"name" gorm:"type:varchar(100);not null;default:'Default Token'"`
	IsActive   bool       `json:"is_active" gorm:"default:true;index"`
	ExpiresAt  *time.Time `json:"expires_at" gorm:"index"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
}

// TableName overrides the table name
func (APIToken) TableName() string {
	return "api_tokens"
}

// HashToken creates a SHA256 hash of the provided token
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// IsExpired checks if the token is expired
func (t *APIToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

// IsValid checks if the token is active and not expired
func (t *APIToken) IsValid() bool {
	return t.IsActive && !t.IsExpired()
}

// UpdateLastUsed updates the last used timestamp
func (t *APIToken) UpdateLastUsed() {
	now := time.Now()
	t.LastUsedAt = &now
}