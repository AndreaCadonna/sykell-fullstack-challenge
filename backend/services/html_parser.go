package services

import (
	"fmt"
	"net/url"
	"strings"
)

// HTMLParser handles HTML analysis and data extraction
type HTMLParser struct{}

// NewHTMLParser creates a new HTML parser instance
func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

// ParsedData contains all extracted information from an HTML page
type ParsedData struct {
	HTMLVersion   *string        `json:"html_version"`   // DOCTYPE analysis
	PageTitle     *string        `json:"page_title"`     // <title> tag content
	HeadingCounts map[string]int `json:"heading_counts"` // h1-h6 counts
	InternalLinks []LinkInfo     `json:"internal_links"` // same domain links
	ExternalLinks []LinkInfo     `json:"external_links"` // external domain links
	HasLoginForm  bool           `json:"has_login_form"` // form with password input
	ParseErrors   []string       `json:"parse_errors"`   // non-fatal parse issues
}

// LinkInfo contains information about a discovered link
type LinkInfo struct {
	URL     string `json:"url"`      // The href URL
	Text    string `json:"text"`     // Link anchor text or alt text
	IsImage bool   `json:"is_image"` // True if this is an image link
}

// Parse analyzes HTML content and extracts all required data
func (p *HTMLParser) Parse(html string, baseURL string) (*ParsedData, error) {
	data := &ParsedData{
		HeadingCounts: make(map[string]int),
		InternalLinks: make([]LinkInfo, 0),
		ExternalLinks: make([]LinkInfo, 0),
		ParseErrors:   make([]string, 0),
		HasLoginForm:  false,
	}

	// Parse base URL for link categorization
	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	baseDomain := strings.ToLower(baseURLParsed.Host)

	// Extract HTML version from DOCTYPE
	if version := p.extractHTMLVersion(html); version != nil {
		data.HTMLVersion = version
	}

	// Extract page title
	if title := p.extractPageTitle(html); title != nil {
		data.PageTitle = title
	}

	// Count heading tags
	data.HeadingCounts = p.extractHeadingCounts(html)

	// Extract and categorize links
	internal, external := p.extractLinks(html, baseDomain)
	data.InternalLinks = internal
	data.ExternalLinks = external

	// Detect login forms
	data.HasLoginForm = p.detectLoginForm(html)

	return data, nil
}

// extractHTMLVersion analyzes DOCTYPE declaration to determine HTML version
func (p *HTMLParser) extractHTMLVersion(html string) *string {
	// Convert to lowercase for case-insensitive matching
	htmlLower := strings.ToLower(html)

	// Look for DOCTYPE declaration
	if strings.Contains(htmlLower, "<!doctype html>") || strings.Contains(htmlLower, "<!doctype html ") {
		version := "HTML5"
		return &version
	}

	if strings.Contains(htmlLower, "html 4.01") {
		if strings.Contains(htmlLower, "strict") {
			version := "HTML4.01 Strict"
			return &version
		} else if strings.Contains(htmlLower, "transitional") {
			version := "HTML4.01 Transitional"
			return &version
		} else if strings.Contains(htmlLower, "frameset") {
			version := "HTML4.01 Frameset"
			return &version
		}
		version := "HTML4.01"
		return &version
	}

	if strings.Contains(htmlLower, "xhtml 1.0") {
		if strings.Contains(htmlLower, "strict") {
			version := "XHTML1.0 Strict"
			return &version
		} else if strings.Contains(htmlLower, "transitional") {
			version := "XHTML1.0 Transitional"
			return &version
		} else if strings.Contains(htmlLower, "frameset") {
			version := "XHTML1.0 Frameset"
			return &version
		}
		version := "XHTML1.0"
		return &version
	}

	if strings.Contains(htmlLower, "xhtml 1.1") {
		version := "XHTML1.1"
		return &version
	}

	// No recognizable DOCTYPE found
	return nil
}

// extractPageTitle extracts the content of the <title> tag
func (p *HTMLParser) extractPageTitle(html string) *string {
	// Simple regex-free implementation for reliability
	// Look for <title> tag (case insensitive)
	htmlLower := strings.ToLower(html)

	titleStart := strings.Index(htmlLower, "<title")
	if titleStart == -1 {
		return nil
	}

	// Find the end of the opening tag
	titleTagEnd := strings.Index(html[titleStart:], ">")
	if titleTagEnd == -1 {
		return nil
	}
	titleTagEnd += titleStart + 1

	// Find the closing tag
	titleEnd := strings.Index(htmlLower[titleTagEnd:], "</title>")
	if titleEnd == -1 {
		return nil
	}
	titleEnd += titleTagEnd

	// Extract title content
	title := strings.TrimSpace(html[titleTagEnd:titleEnd])
	if title == "" {
		return nil
	}

	return &title
}

// extractHeadingCounts counts occurrences of heading tags (h1-h6)
func (p *HTMLParser) extractHeadingCounts(html string) map[string]int {
	counts := map[string]int{
		"h1": 0, "h2": 0, "h3": 0, "h4": 0, "h5": 0, "h6": 0,
	}

	htmlLower := strings.ToLower(html)

	for level := 1; level <= 6; level++ {
		tag := fmt.Sprintf("<h%d", level)
		key := fmt.Sprintf("h%d", level)

		// Count opening tags
		pos := 0
		for {
			index := strings.Index(htmlLower[pos:], tag)
			if index == -1 {
				break
			}
			pos += index + len(tag)

			// Verify it's a complete tag (followed by space, > or other attributes)
			if pos < len(htmlLower) {
				nextChar := htmlLower[pos]
				if nextChar == '>' || nextChar == ' ' || nextChar == '\t' || nextChar == '\n' {
					counts[key]++
				}
			}
		}
	}

	return counts
}

// extractLinks extracts and categorizes links as internal or external
func (p *HTMLParser) extractLinks(html string, baseDomain string) (internal, external []LinkInfo) {
	internal = make([]LinkInfo, 0)
	external = make([]LinkInfo, 0)

	htmlLower := strings.ToLower(html)

	// Simple implementation - look for href attributes
	pos := 0
	for {
		// Find next href attribute
		hrefIndex := strings.Index(htmlLower[pos:], "href=")
		if hrefIndex == -1 {
			break
		}
		hrefIndex += pos

		// Extract the URL value
		linkURL, linkText := p.extractLinkURL(html, hrefIndex)
		if linkURL == "" {
			pos = hrefIndex + 5
			continue
		}

		// Categorize as internal or external
		if p.isInternalLink(linkURL, baseDomain) {
			internal = append(internal, LinkInfo{
				URL:     linkURL,
				Text:    linkText,
				IsImage: false,
			})
		} else {
			external = append(external, LinkInfo{
				URL:     linkURL,
				Text:    linkText,
				IsImage: false,
			})
		}

		pos = hrefIndex + 5
	}

	return internal, external
}

// extractLinkURL extracts URL and text from href attribute
func (p *HTMLParser) extractLinkURL(html string, hrefIndex int) (string, string) {
	// This is a simplified implementation
	// In production, you'd want proper HTML parsing

	// Find the quote character after href=
	start := hrefIndex + 5
	if start >= len(html) {
		return "", ""
	}

	// Skip whitespace
	for start < len(html) && (html[start] == ' ' || html[start] == '\t') {
		start++
	}

	if start >= len(html) {
		return "", ""
	}

	// Check for quote character
	quote := html[start]
	if quote != '"' && quote != '\'' {
		return "", ""
	}

	start++ // Skip opening quote

	// Find closing quote
	end := start
	for end < len(html) && html[end] != quote {
		end++
	}

	if end >= len(html) {
		return "", ""
	}

	linkURL := html[start:end]

	// Extract link text (simplified - just return URL for now)
	linkText := linkURL

	return linkURL, linkText
}

// isInternalLink determines if a link is internal to the base domain
func (p *HTMLParser) isInternalLink(linkURL, baseDomain string) bool {
	// Handle relative URLs
	if strings.HasPrefix(linkURL, "/") ||
		strings.HasPrefix(linkURL, "#") ||
		strings.HasPrefix(linkURL, "?") ||
		(!strings.HasPrefix(linkURL, "http://") && !strings.HasPrefix(linkURL, "https://")) {
		return true
	}

	// Parse absolute URL
	parsedURL, err := url.Parse(linkURL)
	if err != nil {
		return false
	}

	linkDomain := strings.ToLower(parsedURL.Host)
	return linkDomain == baseDomain
}

// detectLoginForm checks if the page contains a login form
func (p *HTMLParser) detectLoginForm(html string) bool {
	htmlLower := strings.ToLower(html)

	// Look for forms containing password inputs
	formStart := 0
	for {
		formIndex := strings.Index(htmlLower[formStart:], "<form")
		if formIndex == -1 {
			break
		}
		formIndex += formStart

		// Find the end of this form
		formEnd := strings.Index(htmlLower[formIndex:], "</form>")
		if formEnd == -1 {
			// No closing form tag found, check rest of document
			formEnd = len(htmlLower)
		} else {
			formEnd += formIndex
		}

		// Check if this form contains a password input
		formContent := htmlLower[formIndex:formEnd]
		if strings.Contains(formContent, "type=\"password\"") ||
			strings.Contains(formContent, "type='password'") ||
			strings.Contains(formContent, "type=password") {
			return true
		}

		formStart = formEnd
	}

	return false
}
