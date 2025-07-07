# Backend Documentation

## Overview

The backend is built with Go 1.22 using the Gin web framework, GORM for database operations, and follows clean architecture principles with clear separation of concerns.

## Architecture

# Add this section to your docs/BACKEND.md file after the existing "Project Structure" section:

### Project Structure

```
backend/
├── main.go                 # Application entry point with CrawlManager
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── .air.toml              # Hot reload configuration
├── Dockerfile             # Multi-stage container build
├── database/
│   └── connection.go      # Database connection and configuration
├── models/
│   ├── url.go            # URL entity model
│   ├── crawl_result.go   # Crawl results model
│   ├── found_link.go     # Found links model
│   └── api_token.go      # API token model
├── dto/
│   ├── requests.go       # Request data transfer objects
│   └── responses.go      # Response data transfer objects
├── handlers/
│   ├── auth.go          # Authentication handlers
│   ├── urls.go          # URL management handlers
│   └── crawl.go         # Crawl control handlers
├── services/             # Business logic layer
│   ├── crawler.go       # HTTP client and safety measures
│   ├── html_parser.go   # HTML analysis and data extraction
│   └── crawl_manager.go # Background job processing
├── middleware/
│   └── auth.go          # Authentication middleware
└── integration_test.go   # Integration tests
```

### Design Patterns

**Repository Pattern**: Models encapsulate database operations  
**DTO Pattern**: Separate request/response objects from domain models  
**Middleware Pattern**: Cross-cutting concerns like authentication  
**Handler Pattern**: HTTP request processing and response formatting  
**Service Pattern**: Business logic separation and dependency injection (NEW)  
**Worker Pattern**: Background job processing with queue management (NEW)  
**Strategy Pattern**: Configurable crawling behavior and error handling (NEW)

## Core Components

### Main Application (`main.go`) - Enhanced

The application entry point now includes:

- Initializes database connection
- **Creates and starts CrawlManager** (NEW)
- **Sets up graceful shutdown handling** (NEW)
- Configures Gin router and middleware
- Sets up CORS for frontend integration
- Defines route groups and handlers (including crawl endpoints)
- Starts the HTTP server

```go
func main() {
    // Database initialization
    database.Connect()

    // CrawlManager initialization
    crawlManager := services.NewCrawlManager()
    crawlManager.Start()
    defer crawlManager.Stop()

    // Graceful shutdown
    setupGracefulShutdown(crawlManager)

    // Router setup
    router := gin.New()
    router.Use(gin.Logger(), gin.Recovery())

    // CORS configuration
    router.Use(cors.New(corsConfig))

    // Handler initialization (including crawl handler)
    crawlHandler := handlers.NewCrawlHandler(crawlManager)

    // Route groups (including crawl routes)
    protected.POST("/urls/:id/crawl", crawlHandler.StartCrawl)
    protected.GET("/urls/:id/crawl/status", crawlHandler.GetCrawlStatus)
    protected.POST("/crawls/bulk", crawlHandler.StartBulkCrawl)
    protected.GET("/crawls/queue/status", crawlHandler.GetQueueStatus)
}
```

### Crawl Handlers (`handlers/crawl.go`)

Handles crawl-related HTTP requests:

- **StartCrawl**: Initiates background crawling for a URL
- **GetCrawlStatus**: Returns real-time crawl status and results
- **StartBulkCrawl**: Manages multiple URL crawling operations
- **GetQueueStatus**: Provides system-wide crawl statistics

```go
type CrawlHandler struct {
    crawlManager *services.CrawlManager
}

func (h *CrawlHandler) StartCrawl(c *gin.Context) {
    // Validate URL exists and is not running
    // Queue URL for background processing
    // Return immediate response (non-blocking)
}
```

### CrawlManager Service (`services/crawl_manager.go`)

Core service that orchestrates the crawling workflow:

- **Background Job Processing**: Goroutine-based queue worker
- **Queue Management**: In-memory job queue with 100-job capacity
- **Database Integration**: Transaction-safe result storage
- **Error Handling**: Comprehensive failure recovery
- **Status Management**: URL status lifecycle management

```go
type CrawlManager struct {
    crawler   *CrawlerService  // HTTP client service
    queue     chan *CrawlJob   // Job queue
    isRunning bool             // Service state
    queueSize int              // Maximum queue capacity
}

func (cm *CrawlManager) processQueue() {
    // Background goroutine processes jobs
    for job := range cm.queue {
        cm.processSingleJob(job)
        time.Sleep(cm.crawler.config.RateLimit)
    }
}
```

### Crawler Service (`services/crawler.go`)

HTTP client with safety measures and content fetching:

- **HTTP Client Configuration**: Timeouts, redirects, size limits
- **Error Classification**: Network, HTTP, content, and parse errors
- **Content Validation**: HTML-only content filtering
- **Safety Measures**: 5MB size limit, 30-second timeout, 5-redirect max

```go
type CrawlerService struct {
    config *CrawlerConfig
    client *http.Client
    parser *HTMLParser
}

type CrawlerConfig struct {
    MaxPageSize    int64         // 5MB limit
    RequestTimeout time.Duration // 30 seconds
    MaxRedirects   int           // 5 redirects max
    UserAgent      string        // "WebCrawler/1.0"
    RateLimit      time.Duration // 1 second between requests
}
```

### HTML Parser (`services/html_parser.go`)

HTML analysis and data extraction engine:

- **HTML Version Detection**: DOCTYPE parsing for HTML5, HTML4.01, XHTML
- **Content Extraction**: Title, headings, links, forms
- **Link Categorization**: Internal vs external link classification
- **Form Detection**: Login form identification
- **Error Resilience**: Continues parsing despite individual failures

```go
type HTMLParser struct{}

type ParsedData struct {
    HTMLVersion   *string           // "HTML5", "HTML4.01 Strict", etc.
    PageTitle     *string           // <title> content
    HeadingCounts map[string]int    // {"h1": 2, "h2": 5, ...}
    InternalLinks []LinkInfo        // Same-domain links
    ExternalLinks []LinkInfo        // External-domain links
    HasLoginForm  bool              // Login form detected
    ParseErrors   []string          // Non-fatal issues
}
```

## Service Layer Architecture

### Dependency Injection Pattern

```go
// Services are injected into handlers
func NewCrawlHandler(crawlManager *services.CrawlManager) *CrawlHandler {
    return &CrawlHandler{
        crawlManager: crawlManager,
    }
}

// CrawlManager depends on CrawlerService
func NewCrawlManager() *CrawlManager {
    return &CrawlManager{
        crawler: NewCrawlerService(nil), // Default config
        queue:   make(chan *CrawlJob, 100),
    }
}
```

### Background Processing Architecture

```go
// Goroutine-based job processing
func (cm *CrawlManager) Start() {
    cm.isRunning = true
    go cm.processQueue() // Background worker
}

// Graceful shutdown
func (cm *CrawlManager) Stop() {
    cm.isRunning = false
    close(cm.queue)
}
```

### Error Handling Strategy

```go
// Layered error handling:
// 1. Network layer (CrawlerService)
// 2. Parse layer (HTMLParser)
// 3. Service layer (CrawlManager)
// 4. Handler layer (CrawlHandler)

func (cm *CrawlManager) handleCrawlFailure(job *CrawlJob, err error, duration time.Duration) {
    log.Printf("Crawl failed for URL ID=%d: %v", job.URLID, err)
    errorMsg := err.Error()
    cm.updateURLStatus(job.URLID, models.StatusError, &errorMsg)
}
```

## Crawl Workflow Implementation

### Complete Crawl Process

```go
// 1. HTTP Request
POST /api/urls/1/crawl
     ↓
// 2. Handler Processing
func (h *CrawlHandler) StartCrawl(c *gin.Context) {
    // Validate URL
    // Queue job
    // Return immediate response
}
     ↓
// 3. Background Processing
func (cm *CrawlManager) processSingleJob(job *CrawlJob) {
    // Update status to "running"
    // Fetch URL content
    // Parse HTML
    // Save results
    // Update status to "completed"
}
     ↓
// 4. Data Storage
Database Transaction:
  - Save crawl_results
  - Save found_links
  - Update URL status
```

### Queue Management

```go
// Non-blocking queue operations
func (cm *CrawlManager) QueueURL(urlID uint, url string) error {
    job := &CrawlJob{URLID: urlID, URL: url, QueuedAt: time.Now()}

    select {
    case cm.queue <- job:
        return nil // Success
    default:
        return fmt.Errorf("queue full") // Fail fast
    }
}
```

### Rate Limiting Implementation

```go
// Respectful crawling with delays
func (cm *CrawlManager) processQueue() {
    for job := range cm.queue {
        cm.processSingleJob(job)
        time.Sleep(cm.crawler.config.RateLimit) // 1 second delay
    }
}
```

## Safety and Performance Features

### Resource Protection

```go
// HTTP client with safety limits
client := &http.Client{
    Timeout: 30 * time.Second,
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        if len(via) >= 5 {
            return errors.New("too many redirects")
        }
        return nil
    },
}

// Content size limiting
limitedReader := io.LimitReader(resp.Body, 5*1024*1024) // 5MB max
```

### Error Recovery

```go
// Transaction-safe database operations
tx := database.DB.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// Save results
if err := saveCrawlResults(tx, data); err != nil {
    tx.Rollback()
    return err
}

tx.Commit() // All or nothing
```

### Monitoring Integration

```go
// Health check includes crawl manager status
router.GET("/health", func(c *gin.Context) {
    response := gin.H{
        "status":        "healthy",
        "crawl_manager": crawlManager.GetQueueStatus(),
    }
    c.JSON(200, response)
})
```

## Performance Characteristics

### Throughput Metrics

- **Queue Capacity**: 100 concurrent jobs
- **Processing Rate**: ~60 URLs/minute (1-second rate limit)
- **Page Size Limit**: 5MB maximum per page
- **Link Storage**: Up to 200 links per page

### Memory Usage

- **Queue Memory**: ~8KB per queued job (100 jobs = ~800KB)
- **Page Content**: Limited to 5MB per active crawl
- **Parser Memory**: Minimal (streaming-based parsing)
- **Database Connections**: Pooled (10 idle, 100 max)
