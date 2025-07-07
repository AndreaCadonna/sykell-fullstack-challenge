package services

import (
	"log"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
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
	URL        string `json:"url"`         // The href URL
	Text       string `json:"text"`        // Link anchor text or alt text
	IsImage    bool   `json:"is_image"`    // True if this is an image link
	IsInternal bool   `json:"is_internal"` // True if internal to domain
}

// Parse analyzes HTML content and extracts all required data
func (p *HTMLParser) Parse(htmlContent string, baseURL string) (*ParsedData, error) {

	log.Printf("DEBUG: Parsing HTML content length: %d, first 200 chars: %s",
		len(htmlContent), htmlContent[:min(200, len(htmlContent))])

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

	// Extract HTML version from DOCTYPE before parsing
	if version := p.extractHTMLVersion(htmlContent); version != nil {
		data.HTMLVersion = version
	}

	// Parse HTML document
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		// Don't fail completely on parse errors, try to extract what we can
		data.ParseErrors = append(data.ParseErrors, "HTML parse error: "+err.Error())
		return data, nil
	}

	// Extract data by traversing the DOM tree
	p.traverseNode(doc, data, baseDomain, baseURL)

	return data, nil
}

// traverseNode recursively walks through the HTML DOM tree
func (p *HTMLParser) traverseNode(n *html.Node, data *ParsedData, baseDomain, baseURL string) {
	if n.Type == html.ElementNode {
		switch strings.ToLower(n.Data) {
		case "title":
			if title := p.extractTextContent(n); title != "" {
				data.PageTitle = &title
			}

		case "h1", "h2", "h3", "h4", "h5", "h6":
			data.HeadingCounts[strings.ToLower(n.Data)]++

		case "a":
			if linkInfo := p.extractLinkInfo(n, baseDomain, baseURL); linkInfo != nil {
				if linkInfo.IsInternal {
					data.InternalLinks = append(data.InternalLinks, *linkInfo)
				} else {
					data.ExternalLinks = append(data.ExternalLinks, *linkInfo)
				}
			}

		case "form":
			if p.hasPasswordInput(n) {
				data.HasLoginForm = true
			}
		}
	}

	// Recursively process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.traverseNode(c, data, baseDomain, baseURL)
	}
}

// extractTextContent extracts the text content from a node and its children
func (p *HTMLParser) extractTextContent(n *html.Node) string {
	var text strings.Builder

	var extractText func(*html.Node)
	extractText = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(n)
	return strings.TrimSpace(text.String())
}

// extractLinkInfo extracts link information from an anchor tag
func (p *HTMLParser) extractLinkInfo(n *html.Node, baseDomain, baseURL string) *LinkInfo {
	var href string

	// Find href attribute
	for _, attr := range n.Attr {
		if strings.ToLower(attr.Key) == "href" {
			href = strings.TrimSpace(attr.Val)
			break
		}
	}

	// Skip empty or invalid hrefs
	if href == "" || href == "#" || strings.HasPrefix(href, "javascript:") || strings.HasPrefix(href, "mailto:") {
		return nil
	}

	// Extract link text
	linkText := p.extractTextContent(n)
	if linkText == "" {
		linkText = href // Fallback to URL if no text
	}

	// Resolve relative URLs
	resolvedURL := p.resolveURL(href, baseURL)
	if resolvedURL == "" {
		return nil
	}

	// Determine if link is internal
	isInternal := p.isInternalLink(resolvedURL, baseDomain)

	return &LinkInfo{
		URL:        resolvedURL,
		Text:       linkText,
		IsImage:    false, // Could be enhanced to detect image links
		IsInternal: isInternal,
	}
}

// resolveURL resolves relative URLs against the base URL
func (p *HTMLParser) resolveURL(href, baseURL string) string {
	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	// Parse href (could be relative or absolute)
	ref, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// Resolve relative URL
	resolved := base.ResolveReference(ref)
	return resolved.String()
}

// isInternalLink determines if a link is internal to the base domain
func (p *HTMLParser) isInternalLink(linkURL, baseDomain string) bool {
	parsedURL, err := url.Parse(linkURL)
	if err != nil {
		return false
	}

	linkDomain := strings.ToLower(parsedURL.Host)
	return linkDomain == baseDomain
}

// hasPasswordInput checks if a form contains a password input
func (p *HTMLParser) hasPasswordInput(formNode *html.Node) bool {
	var hasPassword bool

	var checkInputs func(*html.Node)
	checkInputs = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.ToLower(n.Data) == "input" {
			for _, attr := range n.Attr {
				if strings.ToLower(attr.Key) == "type" && strings.ToLower(attr.Val) == "password" {
					hasPassword = true
					return
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			checkInputs(c)
		}
	}

	checkInputs(formNode)
	return hasPassword
}

// extractHTMLVersion analyzes DOCTYPE declaration to determine HTML version
func (p *HTMLParser) extractHTMLVersion(htmlContent string) *string {
	// Convert to lowercase for case-insensitive matching
	htmlLower := strings.ToLower(htmlContent)

	// Look for DOCTYPE declaration
	doctypeRegex := regexp.MustCompile(`<!doctype\s+html[^>]*>`)
	if doctypeRegex.MatchString(htmlLower) {
		// Check for HTML5
		if strings.Contains(htmlLower, "<!doctype html>") ||
			regexp.MustCompile(`<!doctype\s+html\s*>`).MatchString(htmlLower) {
			version := "HTML5"
			return &version
		}
	}

	// Check for HTML 4.01 variants
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

	// Check for XHTML variants
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
