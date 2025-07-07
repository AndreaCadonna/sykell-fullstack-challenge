package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"web-crawler/database"
	"web-crawler/models"

	"gorm.io/gorm"
)

// CrawlJob represents a crawling job to be processed
type CrawlJob struct {
	URLID    uint      `json:"url_id"`
	URL      string    `json:"url"`
	QueuedAt time.Time `json:"queued_at"`
}

// CrawlManager handles background crawling operations
type CrawlManager struct {
	crawler   *CrawlerService
	queue     chan *CrawlJob
	isRunning bool
	queueSize int
}

// NewCrawlManager creates a new crawl manager instance
func NewCrawlManager() *CrawlManager {
	queueSize := 100 // Reasonable queue size for demo

	return &CrawlManager{
		crawler:   NewCrawlerService(nil), // Use default config
		queue:     make(chan *CrawlJob, queueSize),
		isRunning: false,
		queueSize: queueSize,
	}
}

// Start begins processing crawl jobs in the background
func (cm *CrawlManager) Start() {
	if cm.isRunning {
		log.Println("CrawlManager is already running")
		return
	}

	cm.isRunning = true
	log.Println("Starting CrawlManager background processor")

	go cm.processQueue()
}

// Stop stops the crawl manager (graceful shutdown)
func (cm *CrawlManager) Stop() {
	if !cm.isRunning {
		return
	}

	log.Println("Stopping CrawlManager...")
	cm.isRunning = false
	close(cm.queue)
}

// QueueURL adds a URL to the crawling queue
func (cm *CrawlManager) QueueURL(urlID uint, url string) error {
	if !cm.isRunning {
		return fmt.Errorf("crawl manager is not running")
	}

	job := &CrawlJob{
		URLID:    urlID,
		URL:      url,
		QueuedAt: time.Now(),
	}

	// Non-blocking queue insertion
	select {
	case cm.queue <- job:
		log.Printf("Queued URL for crawling: ID=%d, URL=%s", urlID, url)
		log.Printf("Queue status after queuing: length=%d, size=%d", len(cm.queue), cm.queueSize)
		return nil
	default:
		log.Printf("Queue is full, cannot add URL ID=%d", urlID)
		return fmt.Errorf("crawl queue is full (size: %d)", cm.queueSize)
	}
}

// GetQueueStatus returns information about the current queue state
func (cm *CrawlManager) GetQueueStatus() map[string]interface{} {
	return map[string]interface{}{
		"is_running":   cm.isRunning,
		"queue_length": len(cm.queue),
		"queue_size":   cm.queueSize,
	}
}

// processQueue continuously processes jobs from the queue
func (cm *CrawlManager) processQueue() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("CrawlManager panic recovered: %v", r)
			cm.isRunning = false
		}
	}()

	log.Println("CrawlManager processor started")

	for job := range cm.queue {
		log.Printf("Received job from queue: ID=%d, URL=%s", job.URLID, job.URL)

		if !cm.isRunning {
			log.Println("CrawlManager stopping, remaining jobs will be lost")
			break
		}

		log.Printf("About to process job: ID=%d", job.URLID)
		cm.processSingleJob(job)
		log.Printf("Finished processing job: ID=%d", job.URLID)

		// Rate limiting: wait between jobs
		time.Sleep(cm.crawler.config.RateLimit)
	}

	log.Println("CrawlManager processor stopped")
}

// processSingleJob handles the crawling of a single URL
func (cm *CrawlManager) processSingleJob(job *CrawlJob) {
	log.Printf("Processing crawl job: ID=%d, URL=%s", job.URLID, job.URL)

	// Update URL status to "running"
	if err := cm.updateURLStatus(job.URLID, models.StatusRunning, nil); err != nil {
		log.Printf("Failed to update URL status to running: %v", err)
		return
	}

	// Perform the actual crawl
	startTime := time.Now()
	result, err := cm.performCrawl(job)
	duration := time.Since(startTime)

	if err != nil {
		// Handle crawl failure
		cm.handleCrawlFailure(job, err, duration)
		return
	}

	// Handle crawl success
	cm.handleCrawlSuccess(job, result, duration)
}

// performCrawl executes the actual crawling and parsing
func (cm *CrawlManager) performCrawl(job *CrawlJob) (*ParsedData, error) {
	// Fetch the URL
	response, err := cm.crawler.FetchURL(job.URL)
	if err != nil {
		return nil, err
	}

	// Parse the HTML content
	parsedData, err := cm.crawler.parser.Parse(response.HTML, job.URL)
	if err != nil {
		return nil, err
	}

	return parsedData, nil
}

// handleCrawlSuccess processes successful crawl results
func (cm *CrawlManager) handleCrawlSuccess(job *CrawlJob, data *ParsedData, duration time.Duration) {
	log.Printf("Crawl successful for URL ID=%d, duration=%v", job.URLID, duration)

	// Start database transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction panic for URL ID=%d: %v", job.URLID, r)
		}
	}()

	// Save crawl results
	if err := cm.saveCrawlResults(tx, job.URLID, data, duration); err != nil {
		tx.Rollback()
		log.Printf("Failed to save crawl results for URL ID=%d: %v", job.URLID, err)
		errorMsg := err.Error()
		cm.updateURLStatus(job.URLID, models.StatusError, &errorMsg)
		return
	}

	// Save found links
	if err := cm.saveFoundLinks(tx, job.URLID, data); err != nil {
		tx.Rollback()
		log.Printf("Failed to save found links for URL ID=%d: %v", job.URLID, err)
		errorMsg := err.Error()
		cm.updateURLStatus(job.URLID, models.StatusError, &errorMsg)
		return
	}

	// Update URL status to completed
	if err := cm.updateURLStatusTx(tx, job.URLID, models.StatusCompleted, nil); err != nil {
		tx.Rollback()
		log.Printf("Failed to update URL status to completed for ID=%d: %v", job.URLID, err)
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction for URL ID=%d: %v", job.URLID, err)
		errorMsg := err.Error()
		cm.updateURLStatus(job.URLID, models.StatusError, &errorMsg)
		return
	}

	log.Printf("Crawl completed successfully for URL ID=%d", job.URLID)
}

// handleCrawlFailure processes failed crawl attempts
func (cm *CrawlManager) handleCrawlFailure(job *CrawlJob, err error, duration time.Duration) {
	log.Printf("Crawl failed for URL ID=%d: %v (duration=%v)", job.URLID, err, duration)

	errorMsg := err.Error()
	if updateErr := cm.updateURLStatus(job.URLID, models.StatusError, &errorMsg); updateErr != nil {
		log.Printf("Failed to update URL status to error for ID=%d: %v", job.URLID, updateErr)
	}
}

// updateURLStatus updates the status of a URL in the database
func (cm *CrawlManager) updateURLStatus(urlID uint, status models.URLStatus, errorMsg *string) error {
	return cm.updateURLStatusTx(database.DB, urlID, status, errorMsg)
}

// updateURLStatusTx updates URL status within a transaction
func (cm *CrawlManager) updateURLStatusTx(tx *gorm.DB, urlID uint, status models.URLStatus, errorMsg *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if errorMsg != nil {
		updates["error_message"] = *errorMsg
	} else {
		updates["error_message"] = nil
	}

	return tx.Model(&models.URL{}).Where("id = ?", urlID).Updates(updates).Error
}

// saveCrawlResults saves the parsed HTML data to the crawl_results table
func (cm *CrawlManager) saveCrawlResults(tx *gorm.DB, urlID uint, data *ParsedData, duration time.Duration) error {
	// First, delete any existing crawl result for this URL (re-crawl scenario)
	if err := tx.Where("url_id = ?", urlID).Delete(&models.CrawlResult{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing crawl results: %w", err)
	}

	// Create new crawl result
	durationMs := int(duration.Milliseconds())
	crawlResult := models.CrawlResult{
		URLID:           urlID,
		HTMLVersion:     data.HTMLVersion,
		PageTitle:       data.PageTitle,
		HasLoginForm:    data.HasLoginForm,
		CrawledAt:       time.Now(),
		CrawlDurationMs: &durationMs,
	}

	// Set heading counts
	if h1Count, exists := data.HeadingCounts["h1"]; exists {
		crawlResult.H1Count = h1Count
	}
	if h2Count, exists := data.HeadingCounts["h2"]; exists {
		crawlResult.H2Count = h2Count
	}
	if h3Count, exists := data.HeadingCounts["h3"]; exists {
		crawlResult.H3Count = h3Count
	}
	if h4Count, exists := data.HeadingCounts["h4"]; exists {
		crawlResult.H4Count = h4Count
	}
	if h5Count, exists := data.HeadingCounts["h5"]; exists {
		crawlResult.H5Count = h5Count
	}
	if h6Count, exists := data.HeadingCounts["h6"]; exists {
		crawlResult.H6Count = h6Count
	}

	// Set link counts
	crawlResult.InternalLinksCount = len(data.InternalLinks)
	crawlResult.ExternalLinksCount = len(data.ExternalLinks)

	// For now, set inaccessible links count to 0
	// We'll implement link checking as a future enhancement
	crawlResult.InaccessibleLinksCount = 0

	// Save the crawl result
	if err := tx.Create(&crawlResult).Error; err != nil {
		return fmt.Errorf("failed to create crawl result: %w", err)
	}

	log.Printf("Saved crawl results for URL ID=%d: title=%v, internal_links=%d, external_links=%d",
		urlID,
		formatOptionalString(data.PageTitle),
		len(data.InternalLinks),
		len(data.ExternalLinks))

	return nil
}

// saveFoundLinks saves all discovered links to the found_links table
func (cm *CrawlManager) saveFoundLinks(tx *gorm.DB, urlID uint, data *ParsedData) error {
	// First, delete any existing found links for this URL (re-crawl scenario)
	if err := tx.Where("url_id = ?", urlID).Delete(&models.FoundLink{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing found links: %w", err)
	}

	// Combine internal and external links
	allLinks := make([]models.FoundLink, 0, len(data.InternalLinks)+len(data.ExternalLinks))

	// Process internal links
	for _, link := range data.InternalLinks {
		foundLink := models.FoundLink{
			URLID:      urlID,
			LinkURL:    cm.normalizeURL(link.URL),
			LinkText:   cm.normalizeText(link.Text),
			IsInternal: true,
			CreatedAt:  time.Now(),
		}

		// For now, we don't check accessibility - mark as nil (unchecked)
		// This will be implemented as a future enhancement
		foundLink.IsAccessible = nil
		foundLink.StatusCode = nil
		foundLink.ErrorMessage = nil

		allLinks = append(allLinks, foundLink)
	}

	// Process external links
	for _, link := range data.ExternalLinks {
		foundLink := models.FoundLink{
			URLID:      urlID,
			LinkURL:    cm.normalizeURL(link.URL),
			LinkText:   cm.normalizeText(link.Text),
			IsInternal: false,
			CreatedAt:  time.Now(),
		}

		// For now, we don't check accessibility - mark as nil (unchecked)
		foundLink.IsAccessible = nil
		foundLink.StatusCode = nil
		foundLink.ErrorMessage = nil

		allLinks = append(allLinks, foundLink)
	}

	// Batch insert links (limit to prevent huge inserts)
	const maxLinksToSave = 200 // Reasonable limit for demo
	if len(allLinks) > maxLinksToSave {
		log.Printf("Limiting found links from %d to %d for URL ID=%d", len(allLinks), maxLinksToSave, urlID)
		allLinks = allLinks[:maxLinksToSave]
	}

	// Save links in batches to avoid large transaction
	const batchSize = 50
	for i := 0; i < len(allLinks); i += batchSize {
		end := i + batchSize
		if end > len(allLinks) {
			end = len(allLinks)
		}

		batch := allLinks[i:end]
		if err := tx.Create(&batch).Error; err != nil {
			return fmt.Errorf("failed to create found links batch %d-%d: %w", i, end-1, err)
		}
	}

	log.Printf("Saved %d found links for URL ID=%d (%d internal, %d external)",
		len(allLinks), urlID, len(data.InternalLinks), len(data.ExternalLinks))

	return nil
}

// normalizeURL cleans and validates a URL string
func (cm *CrawlManager) normalizeURL(rawURL string) string {
	// Trim whitespace
	url := strings.TrimSpace(rawURL)

	// Limit URL length to prevent database issues
	const maxURLLength = 2000
	if len(url) > maxURLLength {
		url = url[:maxURLLength]
	}

	return url
}

// normalizeText cleans and validates link text
func (cm *CrawlManager) normalizeText(text string) *string {
	// Trim whitespace
	cleaned := strings.TrimSpace(text)

	// Remove newlines and excessive whitespace
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\r", " ")
	cleaned = strings.ReplaceAll(cleaned, "\t", " ")

	// Collapse multiple spaces
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	// Limit text length
	const maxTextLength = 500
	if len(cleaned) > maxTextLength {
		cleaned = cleaned[:maxTextLength]
	}

	// Return nil for empty text
	if cleaned == "" {
		return nil
	}

	return &cleaned
}

// formatOptionalString formats a string pointer for logging
func formatOptionalString(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("'%s'", *s)
}
