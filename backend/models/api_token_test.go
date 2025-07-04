package models

import (
	"testing"
	"time"
)

func TestHashToken(t *testing.T) {
	token := "test-token-123"
	hash1 := HashToken(token)
	hash2 := HashToken(token)
	
	// Same token should produce same hash
	if hash1 != hash2 {
		t.Errorf("Same token produced different hashes: %s != %s", hash1, hash2)
	}
	
	// Different tokens should produce different hashes
	differentToken := "different-token-456"
	hash3 := HashToken(differentToken)
	
	if hash1 == hash3 {
		t.Errorf("Different tokens produced same hash: %s", hash1)
	}
	
	// Hash should be 64 characters (SHA256 hex)
	if len(hash1) != 64 {
		t.Errorf("Expected hash length to be 64, got %d", len(hash1))
	}
}

func TestAPITokenIsExpired(t *testing.T) {
	// Test token without expiration
	token := APIToken{}
	if token.IsExpired() {
		t.Error("Token without expiration should not be expired")
	}
	
	// Test token with future expiration
	futureTime := time.Now().Add(time.Hour)
	token.ExpiresAt = &futureTime
	if token.IsExpired() {
		t.Error("Token with future expiration should not be expired")
	}
	
	// Test token with past expiration
	pastTime := time.Now().Add(-time.Hour)
	token.ExpiresAt = &pastTime
	if !token.IsExpired() {
		t.Error("Token with past expiration should be expired")
	}
}

func TestAPITokenIsValid(t *testing.T) {
	// Test active, non-expired token
	token := APIToken{
		IsActive: true,
	}
	if !token.IsValid() {
		t.Error("Active, non-expired token should be valid")
	}
	
	// Test inactive token
	token.IsActive = false
	if token.IsValid() {
		t.Error("Inactive token should not be valid")
	}
	
	// Test expired token
	token.IsActive = true
	pastTime := time.Now().Add(-time.Hour)
	token.ExpiresAt = &pastTime
	if token.IsValid() {
		t.Error("Expired token should not be valid")
	}
}

func TestAPITokenUpdateLastUsed(t *testing.T) {
	token := APIToken{}
	
	// Initially should be nil
	if token.LastUsedAt != nil {
		t.Error("LastUsedAt should initially be nil")
	}
	
	// After update should be set
	token.UpdateLastUsed()
	if token.LastUsedAt == nil {
		t.Error("LastUsedAt should be set after update")
	}
	
	// Should be recent (within last second)
	timeDiff := time.Since(*token.LastUsedAt)
	if timeDiff > time.Second {
		t.Errorf("LastUsedAt should be recent, but was %v ago", timeDiff)
	}
}

func TestFoundLinkIsBroken(t *testing.T) {
	// Test link without status code
	link := FoundLink{}
	if link.IsBroken() {
		t.Error("Link without status code should not be considered broken")
	}
	
	// Test successful status codes
	successCodes := []int{200, 201, 204, 301, 302}
	for _, code := range successCodes {
		link.StatusCode = &code
		if link.IsBroken() {
			t.Errorf("Link with status code %d should not be considered broken", code)
		}
	}
	
	// Test broken status codes
	brokenCodes := []int{400, 404, 500, 503}
	for _, code := range brokenCodes {
		link.StatusCode = &code
		if !link.IsBroken() {
			t.Errorf("Link with status code %d should be considered broken", code)
		}
	}
}

func TestFoundLinkGetStatusCategory(t *testing.T) {
	link := FoundLink{}
	
	// Test unchecked status
	if link.GetStatusCategory() != "unchecked" {
		t.Error("Link without status code should have 'unchecked' category")
	}
	
	// Test status categories
	testCases := []struct {
		statusCode int
		expected   string
	}{
		{200, "success"},
		{201, "success"},
		{301, "redirect"},
		{302, "redirect"},
		{400, "client_error"},
		{404, "client_error"},
		{500, "server_error"},
		{503, "server_error"},
		{100, "unknown"},
	}
	
	for _, tc := range testCases {
		link.StatusCode = &tc.statusCode
		category := link.GetStatusCategory()
		if category != tc.expected {
			t.Errorf("Status code %d should have category '%s', got '%s'", 
				tc.statusCode, tc.expected, category)
		}
	}
}