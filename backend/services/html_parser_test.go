package services

import (
	"testing"
)

func TestHTMLParser_Parse(t *testing.T) {
	parser := NewHTMLParser()

	testHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page Title</title>
</head>
<body>
    <h1>Main Heading</h1>
    <h2>Sub Heading 1</h2>
    <h2>Sub Heading 2</h2>
    <h3>Sub Sub Heading</h3>
    
    <p>This is a paragraph with <a href="/internal">internal link</a> and 
       <a href="https://external.com">external link</a>.</p>
    
    <a href="https://example.com/page">Another internal link</a>
    <a href="https://google.com">Google</a>
    
    <form>
        <input type="text" name="username" />
        <input type="password" name="password" />
        <button type="submit">Login</button>
    </form>
</body>
</html>
`

	baseURL := "https://example.com"

	result, err := parser.Parse(testHTML, baseURL)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Test HTML version detection
	if result.HTMLVersion == nil || *result.HTMLVersion != "HTML5" {
		t.Errorf("Expected HTML5, got %v", result.HTMLVersion)
	}

	// Test title extraction
	if result.PageTitle == nil || *result.PageTitle != "Test Page Title" {
		t.Errorf("Expected 'Test Page Title', got %v", result.PageTitle)
	}

	// Test heading counts
	expectedHeadings := map[string]int{
		"h1": 1,
		"h2": 2,
		"h3": 1,
		"h4": 0,
		"h5": 0,
		"h6": 0,
	}

	for tag, expectedCount := range expectedHeadings {
		if result.HeadingCounts[tag] != expectedCount {
			t.Errorf("Expected %s count to be %d, got %d", tag, expectedCount, result.HeadingCounts[tag])
		}
	}

	// Test internal links
	if len(result.InternalLinks) != 2 {
		t.Errorf("Expected 2 internal links, got %d", len(result.InternalLinks))
	}

	// Test external links
	if len(result.ExternalLinks) != 2 {
		t.Errorf("Expected 2 external links, got %d", len(result.ExternalLinks))
	}

	// Test login form detection
	if !result.HasLoginForm {
		t.Error("Expected login form to be detected")
	}

	// Print results for debugging
	t.Logf("Parsed data: %+v", result)
	t.Logf("Internal links: %+v", result.InternalLinks)
	t.Logf("External links: %+v", result.ExternalLinks)
}

func TestHTMLParser_ParseSimple(t *testing.T) {
	parser := NewHTMLParser()

	// Simple HTML like example.com
	simpleHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Example Domain</title>
</head>
<body>
    <div>
        <h1>Example Domain</h1>
        <p>This domain is for use in illustrative examples in documents.</p>
        <p><a href="https://www.iana.org/domains/example">More information...</a></p>
    </div>
</body>
</html>
`

	baseURL := "https://example.com"

	result, err := parser.Parse(simpleHTML, baseURL)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify basic extraction works
	if result.PageTitle == nil {
		t.Error("Title should be extracted")
	} else if *result.PageTitle != "Example Domain" {
		t.Errorf("Expected 'Example Domain', got '%s'", *result.PageTitle)
	}

	if result.HeadingCounts["h1"] != 1 {
		t.Errorf("Expected 1 H1 tag, got %d", result.HeadingCounts["h1"])
	}

	if len(result.ExternalLinks) == 0 {
		t.Error("Should find at least one external link")
	}

	t.Logf("Simple parse result: Title=%v, H1=%d, External Links=%d",
		result.PageTitle, result.HeadingCounts["h1"], len(result.ExternalLinks))
}

func TestHTMLParser_EdgeCases(t *testing.T) {
	parser := NewHTMLParser()

	// Test malformed HTML
	malformedHTML := `
<html>
<title>Broken Title
<h1>Unclosed heading
<a href="/test">Link without closing tag
<form>
<input type="password"
</form>
`

	result, err := parser.Parse(malformedHTML, "https://test.com")
	if err != nil {
		t.Fatalf("Should handle malformed HTML gracefully: %v", err)
	}

	// Should still extract some data despite malformed HTML
	t.Logf("Malformed HTML result: %+v", result)

	// Test empty HTML
	emptyResult, err := parser.Parse("", "https://test.com")
	if err != nil {
		t.Fatalf("Should handle empty HTML: %v", err)
	}

	if emptyResult.PageTitle != nil {
		t.Error("Empty HTML should not have title")
	}
}
