package dto

import (
	"testing"
)

func TestAddURLRequestValidate(t *testing.T) {
	testCases := []struct {
		name      string
		url       string
		shouldErr bool
	}{
		{
			name:      "valid HTTP URL",
			url:       "http://example.com",
			shouldErr: false,
		},
		{
			name:      "valid HTTPS URL",
			url:       "https://example.com/path",
			shouldErr: false,
		},
		{
			name:      "empty URL",
			url:       "",
			shouldErr: true,
		},
		{
			name:      "whitespace only URL",
			url:       "   ",
			shouldErr: true,
		},
		{
			name:      "URL without protocol",
			url:       "example.com",
			shouldErr: true,
		},
		{
			name:      "URL with invalid protocol",
			url:       "ftp://example.com",
			shouldErr: true,
		},
		{
			name:      "URL without host",
			url:       "http://",
			shouldErr: true,
		},
		{
			name:      "malformed URL",
			url:       "not-a-url",
			shouldErr: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := AddURLRequest{URL: tc.url}
			err := req.Validate()
			
			if tc.shouldErr && err == nil {
				t.Errorf("Expected validation to fail for URL: %s", tc.url)
			}
			
			if !tc.shouldErr && err != nil {
				t.Errorf("Expected validation to pass for URL: %s, got error: %v", tc.url, err)
			}
		})
	}
}

func TestAddURLRequestNormalize(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trim whitespace",
			input:    "  https://example.com  ",
			expected: "https://example.com",
		},
		{
			name:     "lowercase scheme and host",
			input:    "HTTPS://EXAMPLE.COM/Path",
			expected: "https://example.com/Path",
		},
		{
			name:     "remove trailing slash from path",
			input:    "https://example.com/path/",
			expected: "https://example.com/path",
		},
		{
			name:     "keep root path slash",
			input:    "https://example.com/",
			expected: "https://example.com/",
		},
		{
			name:     "preserve query parameters",
			input:    "https://example.com/path?param=value",
			expected: "https://example.com/path?param=value",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := AddURLRequest{URL: tc.input}
			req.Normalize()
			
			if req.URL != tc.expected {
				t.Errorf("Expected normalized URL to be '%s', got '%s'", tc.expected, req.URL)
			}
		})
	}
}

func TestPaginationRequestValidate(t *testing.T) {
	// Test valid pagination request
	req := PaginationRequest{
		Page:     1,
		PageSize: 20,
		SortBy:   "created_at",
		SortDir:  "desc",
		Status:   "queued",
	}
	
	if err := req.Validate(); err != nil {
		t.Errorf("Valid pagination request should not return error: %v", err)
	}
	
	// Test invalid sort direction
	req.SortDir = "invalid"
	if err := req.Validate(); err == nil {
		t.Error("Invalid sort direction should return error")
	}
	
	// Test invalid sort field
	req.SortDir = "desc"
	req.SortBy = "invalid_field"
	if err := req.Validate(); err == nil {
		t.Error("Invalid sort field should return error")
	}
	
	// Test invalid status
	req.SortBy = "created_at"
	req.Status = "invalid_status"
	if err := req.Validate(); err == nil {
		t.Error("Invalid status should return error")
	}
}

func TestPaginationRequestGetOffset(t *testing.T) {
	testCases := []struct {
		page     int
		pageSize int
		expected int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 10, 20},
		{5, 25, 100},
	}
	
	for _, tc := range testCases {
		req := PaginationRequest{
			Page:     tc.page,
			PageSize: tc.pageSize,
		}
		
		offset := req.GetOffset()
		if offset != tc.expected {
			t.Errorf("Expected offset %d for page %d and page size %d, got %d",
				tc.expected, tc.page, tc.pageSize, offset)
		}
	}
}

func TestPaginationRequestGetOrderClause(t *testing.T) {
	req := PaginationRequest{
		SortBy:  "created_at",
		SortDir: "desc",
	}
	
	expected := "created_at DESC"
	orderClause := req.GetOrderClause()
	
	if orderClause != expected {
		t.Errorf("Expected order clause '%s', got '%s'", expected, orderClause)
	}
	
	// Test ascending order
	req.SortDir = "asc"
	expected = "created_at ASC"
	orderClause = req.GetOrderClause()
	
	if orderClause != expected {
		t.Errorf("Expected order clause '%s', got '%s'", expected, orderClause)
	}
}