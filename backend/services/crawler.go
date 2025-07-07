package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// CrawlerConfig holds configuration for the crawler service
type CrawlerConfig struct {
	MaxPageSize    int64         // Maximum page size in bytes (5MB)
	RequestTimeout time.Duration // HTTP request timeout (30 seconds)
	MaxRedirects   int           // Maximum number of redirects to follow (5)
	UserAgent      string        // User agent string
	RateLimit      time.Duration // Delay between requests (1 second)
}

// DefaultCrawlerConfig returns a safe default configuration
func DefaultCrawlerConfig() *CrawlerConfig {
	return &CrawlerConfig{
		MaxPageSize:    5 * 1024 * 1024, // 5MB
		RequestTimeout: 30 * time.Second,
		MaxRedirects:   5,
		UserAgent:      "WebCrawler/1.0 (+https://github.com/your-repo/web-crawler)",
		RateLimit:      1 * time.Second,
	}
}

// CrawlerService handles web crawling operations
type CrawlerService struct {
	config *CrawlerConfig
	client *http.Client
	parser *HTMLParser
}

// NewCrawlerService creates a new crawler service with the given configuration
func NewCrawlerService(config *CrawlerConfig) *CrawlerService {
	if config == nil {
		config = DefaultCrawlerConfig()
	}

	// Create HTTP client with proper configuration
	client := &http.Client{
		Timeout: config.RequestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent infinite redirect loops
			if len(via) >= config.MaxRedirects {
				return errors.New("too many redirects")
			}
			return nil
		},
	}

	return &CrawlerService{
		config: config,
		client: client,
		parser: NewHTMLParser(), // We'll implement this next
	}
}

// CrawlResponse contains the result of fetching a URL
type CrawlResponse struct {
	HTML         string        // Raw HTML content
	StatusCode   int           // HTTP status code
	ContentType  string        // Content-Type header
	ResponseSize int64         // Size of response in bytes
	Duration     time.Duration // Time taken to fetch
	URL          string        // Final URL (after redirects)
}

// CrawlError represents a crawling error with context
type CrawlError struct {
	Type    string // Error type: "network", "timeout", "too_large", "invalid_url"
	Message string // Human-readable error message
	URL     string // URL that caused the error
	Err     error  // Underlying error
}

func (e *CrawlError) Error() string {
	return e.Message
}

// NewCrawlError creates a new crawl error with the given details
func NewCrawlError(errorType, message, url string, err error) *CrawlError {
	return &CrawlError{
		Type:    errorType,
		Message: message,
		URL:     url,
		Err:     err,
	}
}

// FetchURL fetches and validates a URL with all safety measures
func (c *CrawlerService) FetchURL(rawURL string) (*CrawlResponse, error) {
	startTime := time.Now()

	// Validate URL format
	if err := c.validateURL(rawURL); err != nil {
		return nil, NewCrawlError("invalid_url", "Invalid URL format", rawURL, err)
	}

	// Create HTTP request with proper headers
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, NewCrawlError("invalid_url", "Failed to create HTTP request", rawURL, err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")

	// Perform HTTP request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.classifyNetworkError(rawURL, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 400 {
		return nil, NewCrawlError("http_error",
			fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			rawURL, nil)
	}

	// Validate content type
	contentType := resp.Header.Get("Content-Type")
	if !c.isHTMLContent(contentType) {
		return nil, NewCrawlError("invalid_content",
			fmt.Sprintf("Content-Type '%s' is not HTML", contentType),
			rawURL, nil)
	}

	// Read response body with size limit
	body, err := c.readResponseBody(resp, rawURL)
	if err != nil {
		return nil, err // Already wrapped in CrawlError
	}

	duration := time.Since(startTime)

	return &CrawlResponse{
		HTML:         string(body),
		StatusCode:   resp.StatusCode,
		ContentType:  contentType,
		ResponseSize: int64(len(body)),
		Duration:     duration,
		URL:          resp.Request.URL.String(), // Final URL after redirects
	}, nil
}

// validateURL checks if the URL is valid and uses HTTP/HTTPS
func (c *CrawlerService) validateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must use HTTP or HTTPS protocol")
	}

	if parsedURL.Host == "" {
		return errors.New("URL must have a valid host")
	}

	return nil
}

// classifyNetworkError categorizes network errors for better error handling
func (c *CrawlerService) classifyNetworkError(url string, err error) *CrawlError {
	if errors.Is(err, context.DeadlineExceeded) {
		return NewCrawlError("timeout", "Request timed out", url, err)
	}

	if strings.Contains(err.Error(), "no such host") {
		return NewCrawlError("dns_error", "DNS lookup failed", url, err)
	}

	if strings.Contains(err.Error(), "connection refused") {
		return NewCrawlError("connection_error", "Connection refused", url, err)
	}

	if strings.Contains(err.Error(), "too many redirects") {
		return NewCrawlError("redirect_error", "Too many redirects", url, err)
	}

	// Generic network error
	return NewCrawlError("network", "Network error occurred", url, err)
}

// isHTMLContent checks if the content type indicates HTML content
func (c *CrawlerService) isHTMLContent(contentType string) bool {
	if contentType == "" {
		return true // Assume HTML if no content type
	}

	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "text/html") ||
		strings.Contains(contentType, "application/xhtml+xml")
}

// readResponseBody reads the response body with size limits
func (c *CrawlerService) readResponseBody(resp *http.Response, url string) ([]byte, error) {
	// Create a limited reader to prevent reading too much data
	limitedReader := io.LimitReader(resp.Body, c.config.MaxPageSize+1)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, NewCrawlError("read_error", "Failed to read response body", url, err)
	}

	// Check if we hit the size limit
	if int64(len(body)) > c.config.MaxPageSize {
		return nil, NewCrawlError("too_large",
			fmt.Sprintf("Response size exceeds limit of %d bytes", c.config.MaxPageSize),
			url, nil)
	}

	return body, nil
}
