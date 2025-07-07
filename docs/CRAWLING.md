# Crawling Service Documentation

## Overview

The Web Crawler implements a sophisticated background job processing system that safely crawls websites and extracts structured data. The service is built with Go and follows production-ready patterns for reliability, performance, and error handling.

## Architecture

### Core Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Handler  │    │  CrawlManager   │    │ CrawlerService  │
│                 │    │                 │    │                 │
│ • StartCrawl    │───▶│ • Job Queue     │───▶│ • HTTP Client   │
│ • GetStatus     │    │ • Background    │    │ • HTML Parser   │
│ • BulkCrawl     │    │   Processing    │    │ • Error Handle  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Database     │    │   Job Queue     │    │  Target Website │
│                 │    │                 │    │                 │
│ • URLs          │    │ • In-Memory     │    │ • HTML Content  │
│ • CrawlResults  │    │ • 100 Jobs Max  │    │ • Links         │
│ • FoundLinks    │    │ • Rate Limited  │    │ • Forms         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Service Layers

1. **API Layer** (`handlers/crawl.go`) - HTTP endpoints and request validation
2. **Service Layer** (`services/crawl_manager.go`) - Business logic and orchestration
3. **Worker Layer** (`services/crawler.go`) - HTTP client and content fetching
4. **Parser Layer** (`services/html_parser.go`) - HTML analysis and data extraction
5. **Data Layer** (`models/*.go`) - Database entities and persistence

## Crawl Workflow

### Step-by-Step Process

#### 1. Crawl Initiation

```
POST /api/urls/1/crawl
     ↓
[Handler] Validate URL exists and not running
     ↓
[CrawlManager] Add job to in-memory queue
     ↓
[Response] Immediate 202 Accepted (non-blocking)
```

#### 2. Background Processing

```
[Background Goroutine] Pick up job from queue
     ↓
[Database] Update URL status: queued → running
     ↓
[CrawlerService] Fetch URL content (HTTP request)
     ↓
[HTMLParser] Extract data from HTML
     ↓
[Database Transaction] Save results + update status
```

#### 3. Data Extraction & Storage

```
HTML Content
     ↓
[Parser] Extract:
  • HTML Version (DOCTYPE)
  • Page Title (<title>)
  • Heading Counts (H1-H6)
  • Links (internal vs external)
  • Login Forms (password inputs)
     ↓
[Database] Save to:
  • crawl_results (extracted data)
  • found_links (discovered links)
```

#### 4. Status Updates

```
URL Status Flow:
queued → running → completed/error
     ↑         ↑           ↑
  [Start]  [Process]   [Finish]
```

## Configuration & Safety

### HTTP Client Safety

```go
CrawlerConfig {
    MaxPageSize:    5MB           // Prevent memory issues
    RequestTimeout: 30 seconds    // Prevent hanging requests
    MaxRedirects:   5             // Prevent infinite loops
    UserAgent:      "WebCrawler/1.0"
    RateLimit:      1 second      // Respectful crawling
}
```

### Error Classification

- **Network Errors**: DNS lookup failed, connection refused, timeout
- **HTTP Errors**: 4xx client errors, 5xx server errors
- **Content Errors**: Non-HTML content, page too large
- **Parse Errors**: Malformed HTML, extraction failures

### Resource Limits

- **Queue Size**: 100 jobs maximum
- **Page Size**: 5MB maximum per page
- **Link Limit**: 200 links saved per page
- **Batch Size**: 50 links per database transaction

## Data Extraction Capabilities

### HTML Analysis

#### HTML Version Detection

```go
// Supports detection of:
"HTML5"                 // <!DOCTYPE html>
"HTML4.01 Strict"      // <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
"HTML4.01 Transitional" // <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN"
"XHTML1.0 Strict"      // <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
// ... and other variants
```

#### Link Analysis

```go
type LinkInfo struct {
    URL        string  // href attribute value
    Text       string  // anchor text or alt text
    IsInternal bool    // same domain = internal
}

// Categorization Logic:
• Same domain (example.com) = Internal
• Different domain = External
• Relative URLs (/path, #anchor) = Internal
• Protocol-less URLs = Internal
```

#### Form Detection

```go
// Detects login forms by finding:
<form>
  <input type="password" />
</form>

// Returns: has_login_form = true
```

### Performance Metrics

- **Crawl Duration**: Measured in milliseconds
- **Response Size**: Tracked for monitoring
- **Link Counts**: Internal vs external statistics
- **Error Rates**: Tracked per error type

## Database Schema Integration

### Core Tables

#### URLs Management

```sql
urls:
  id, url, status, error_message, created_at, updated_at

Status Values: 'queued', 'running', 'completed', 'error'
```

#### Crawl Results

```sql
crawl_results:
  id, url_id, html_version, page_title,
  h1_count, h2_count, h3_count, h4_count, h5_count, h6_count,
  internal_links_count, external_links_count, inaccessible_links_count,
  has_login_form, crawled_at, crawl_duration_ms
```

#### Found Links

```sql
found_links:
  id, url_id, link_url, link_text, is_internal,
  is_accessible, status_code, error_message, created_at

Note: is_accessible and status_code are prepared for future link checking
```

### Relationship Design

```
urls (1) ←→ (1) crawl_results    # One crawl result per URL
urls (1) ←→ (∞) found_links      # Many links per URL
```

## Error Handling & Recovery

### Graceful Failure Patterns

#### Partial Data Recovery

```go
// If title extraction succeeds but link parsing fails:
// → Save title, set link counts to 0, mark as completed
// → Better than total failure
```

#### Transaction Safety

```go
// All database operations are wrapped in transactions:
tx.Begin()
  saveResults()
  saveLinks()
  updateStatus()
tx.Commit() // or Rollback on any failure
```

#### Queue Resilience

```go
// Non-blocking queue operations:
select {
case queue <- job:
    // Success
default:
    // Queue full, return error immediately
}
```

### Error Recovery Strategies

1. **Network Timeouts**: Classified and logged, don't crash service
2. **Parse Failures**: Save what was successfully extracted
3. **Database Failures**: Rollback transactions, maintain data integrity
4. **Queue Overflow**: Reject new jobs gracefully, don't block service

## Performance & Scalability

### Current Design Choices

#### Single-Threaded Processing

- **Pros**: Simple, predictable, easy to debug
- **Cons**: Limited throughput
- **Future**: Can scale to worker pools

#### In-Memory Queue

- **Pros**: Fast, no external dependencies
- **Cons**: Jobs lost on restart
- **Future**: Can upgrade to Redis/RabbitMQ

#### Rate Limiting

- **Current**: 1 second between requests
- **Rationale**: Respectful to target websites
- **Future**: Per-domain rate limiting

### Scaling Strategies

#### Horizontal Scaling

```go
// Future: Multiple worker instances
CrawlManager {
    workers: []Worker{
        Worker{id: 1, queue: chan1},
        Worker{id: 2, queue: chan2},
        Worker{id: 3, queue: chan3},
    }
}
```

#### Queue Scaling

```go
// Future: External queue system
type QueueInterface interface {
    Push(job *CrawlJob) error
    Pop() (*CrawlJob, error)
    Size() int
}
```

## Monitoring & Observability

### Health Metrics

```json
GET /health:
{
  "crawl_manager": {
    "is_running": true,
    "queue_length": 2,
    "queue_size": 100
  }
}
```

### Queue Statistics

```json
GET /api/crawls/queue/status:
{
  "queue_manager": {...},
  "database_stats": {
    "queued_count": 5,
    "running_count": 1,
    "completed_count": 23,
    "error_count": 2
  }
}
```

### Detailed Logging

```
[CRAWL] StartCrawl endpoint called
[CRAWL] Found URL: ID=1, URL=https://example.com, Status=queued
[QUEUE] Successfully queued URL: ID=1, URL=https://example.com
[PROCESSOR] Received job from queue: ID=1, URL=https://example.com
[FETCH] Successfully fetched URL: https://example.com, Size: 1547 bytes
[PARSE] Successfully parsed HTML: Title='Example Domain', Links: 1 internal + 1 external
[JOB] Crawl succeeded, handling success: ID=1
```

## Configuration Options

### Environment Variables

```bash
# Rate limiting
CRAWL_RATE_LIMIT=1s          # Time between requests

# Safety limits
MAX_PAGE_SIZE=5242880        # 5MB in bytes
REQUEST_TIMEOUT=30s          # HTTP timeout
MAX_REDIRECTS=5              # Redirect limit

# Queue settings
QUEUE_SIZE=100               # Maximum queued jobs
MAX_LINKS_PER_PAGE=200       # Link storage limit

# Processing
ENABLE_LINK_CHECKING=false   # Future feature flag
```

### Runtime Configuration

```go
// Configurable via code:
config := &CrawlerConfig{
    MaxPageSize:    5 * 1024 * 1024,
    RequestTimeout: 30 * time.Second,
    MaxRedirects:   5,
    UserAgent:      "WebCrawler/1.0",
    RateLimit:      1 * time.Second,
}
```

## Testing Strategy

### Unit Tests

```bash
# Test individual components:
go test ./services/crawler_test.go      # HTTP client
go test ./services/html_parser_test.go  # HTML parsing
go test ./handlers/crawl_test.go        # API endpoints
```

### Integration Tests

```bash
# Test complete workflow:
docker-compose -f docker-compose.test.yml up --build
```

### Manual Testing

```bash
# Complete crawl test:
curl -X POST .../urls -d '{"url":"https://example.com"}'
curl -X POST .../urls/1/crawl
curl .../urls/1/crawl/status
curl .../urls/1/details
```

## Future Enhancements

### Link Accessibility Checking

- **Infrastructure**: Already in database (is_accessible, status_code)
- **Implementation**: Background service to check found links
- **Benefits**: Identify broken links, SEO analysis

### Advanced Scheduling

- **Priority Queues**: High/medium/low priority crawls
- **Scheduled Crawls**: Recurring crawl schedules
- **Dependency Management**: Crawl ordering based on relationships

### Content Analysis

- **Image Analysis**: Extract and analyze images
- **Content Classification**: Categorize page content
- **SEO Analysis**: Meta tags, structured data
- **Performance Metrics**: Page load times, resource analysis

### Distributed Processing

- **Worker Pools**: Multiple concurrent crawlers
- **Microservices**: Separate parsing service
- **Message Queues**: Redis/RabbitMQ integration
- **Load Balancing**: Distribute crawl jobs across instances

## Troubleshooting Guide

### Common Issues

#### Queue Not Processing

```bash
# Check if CrawlManager started:
docker logs backend | grep "CrawlManager processor started"

# Check queue status:
curl .../api/crawls/queue/status
```

#### Crawls Timing Out

```bash
# Check for slow websites:
docker logs backend | grep "FETCH.*ERROR"

# Adjust timeout in config
```

#### Memory Issues

```bash
# Check for large pages:
docker logs backend | grep "too_large"

# Monitor container memory:
docker stats
```

#### Database Connection Issues

```bash
# Check database logs:
docker logs database

# Test connection:
curl .../health
```

### Debug Commands

```bash
# Watch crawl logs in real-time:
docker logs backend -f | grep -E "\[CRAWL\]|\[QUEUE\]|\[PROCESSOR\]"

# Check database status:
docker exec -it database mysql -u crawler_user -pcrawler_password -e "
  SELECT status, COUNT(*) FROM crawler_db.urls GROUP BY status;
"

# Monitor queue length:
watch -n 1 "curl -s .../api/crawls/queue/status | jq '.data.queue_manager.queue_length'"
```

This crawling service provides a robust, production-ready foundation for web crawling operations with proper error handling, monitoring, and scalability considerations.
