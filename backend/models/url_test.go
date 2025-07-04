package models

import (
	"testing"
)

func TestURLStatusConstants(t *testing.T) {
	// Test that URL status constants are defined correctly
	if StatusQueued != "queued" {
		t.Errorf("Expected StatusQueued to be 'queued', got %s", StatusQueued)
	}
	
	if StatusRunning != "running" {
		t.Errorf("Expected StatusRunning to be 'running', got %s", StatusRunning)
	}
	
	if StatusCompleted != "completed" {
		t.Errorf("Expected StatusCompleted to be 'completed', got %s", StatusCompleted)
	}
	
	if StatusError != "error" {
		t.Errorf("Expected StatusError to be 'error', got %s", StatusError)
	}
}

func TestURLTableName(t *testing.T) {
	url := URL{}
	expectedTableName := "urls"
	
	if url.TableName() != expectedTableName {
		t.Errorf("Expected table name to be '%s', got '%s'", expectedTableName, url.TableName())
	}
}

func TestCrawlResultGetHeadingCounts(t *testing.T) {
	result := CrawlResult{
		H1Count: 2,
		H2Count: 5,
		H3Count: 3,
		H4Count: 1,
		H5Count: 0,
		H6Count: 1,
	}
	
	headingCounts := result.GetHeadingCounts()
	
	expected := map[string]int{
		"h1": 2,
		"h2": 5,
		"h3": 3,
		"h4": 1,
		"h5": 0,
		"h6": 1,
	}
	
	for tag, expectedCount := range expected {
		if headingCounts[tag] != expectedCount {
			t.Errorf("Expected %s count to be %d, got %d", tag, expectedCount, headingCounts[tag])
		}
	}
}

func TestCrawlResultGetTotalLinks(t *testing.T) {
	result := CrawlResult{
		InternalLinksCount: 15,
		ExternalLinksCount: 8,
	}
	
	totalLinks := result.GetTotalLinks()
	expected := 23
	
	if totalLinks != expected {
		t.Errorf("Expected total links to be %d, got %d", expected, totalLinks)
	}
}