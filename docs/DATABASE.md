# Database Documentation

## Overview

The Web Crawler uses MySQL 8.0 with GORM as the ORM layer. The database schema is designed to support web crawling operations with proper normalization and indexing for performance.

## Schema Design

### Entity Relationship Diagram

```
urls (1) ←→ (1) crawl_results
urls (1) ←→ (∞) found_links
api_tokens (standalone)
```

## Tables

### 1. `urls` - Target URLs for crawling

| Column | Type | Description | Constraints |
|--------|------|-------------|-------------|
| `id` | BIGINT | Primary key | AUTO_INCREMENT |
| `url` | VARCHAR(2048) | Target URL | NOT NULL, UNIQUE(255) |
| `status` | ENUM | Crawl status | 'queued', 'running', 'completed', 'error' |
| `error_message` | TEXT | Error details | NULLABLE |
| `created_at` | TIMESTAMP | Creation time | DEFAULT CURRENT_TIMESTAMP |
| `updated_at` | TIMESTAMP | Last update | ON UPDATE CURRENT_TIMESTAMP |

**Indexes:**
- Primary: `id`
- Unique: `url(255)` (prevents duplicates)
- Index: `status` (for filtering)
- Index: `created_at` (for sorting)

### 2. `crawl_results` - Extracted data from crawling

| Column | Type | Description | Default |
|--------|------|-------------|---------|
| `id` | BIGINT | Primary key | AUTO_INCREMENT |
| `url_id` | BIGINT | Foreign key to urls | NOT NULL |
| `html_version` | VARCHAR(50) | HTML version (e.g., "HTML5") | NULL |
| `page_title` | VARCHAR(500) | Page title tag | NULL |
| `h1_count` | INT | Number of H1 tags | 0 |
| `h2_count` | INT | Number of H2 tags | 0 |
| `h3_count` | INT | Number of H3 tags | 0 |
| `h4_count` | INT | Number of H4 tags | 0 |
| `h5_count` | INT | Number of H5 tags | 0 |
| `h6_count` | INT | Number of H6 tags | 0 |
| `internal_links_count` | INT | Number of internal links | 0 |
| `external_links_count` | INT | Number of external links | 0 |
| `inaccessible_links_count` | INT | Number of broken links | 0 |
| `has_login_form` | BOOLEAN | Login form detected | FALSE |
| `crawled_at` | TIMESTAMP | Crawl completion time | Set by code |
| `crawl_duration_ms` | INT | Crawl duration in milliseconds | NULL |

**Indexes:**
- Primary: `id`
- Foreign key: `url_id` → `urls(id)` (CASCADE DELETE)
- Index: `crawled_at` (for sorting)

### 3. `found_links` - Links discovered during crawling

| Column | Type | Description | Constraints |
|--------|------|-------------|-------------|
| `id` | BIGINT | Primary key | AUTO_INCREMENT |
| `url_id` | BIGINT | Foreign key to urls | NOT NULL |
| `link_url` | VARCHAR(2048) | Discovered link URL | NOT NULL |
| `link_text` | VARCHAR(500) | Link anchor text | NULLABLE |
| `is_internal` | BOOLEAN | Internal to domain | NOT NULL |
| `is_accessible` | BOOLEAN | Link accessibility status | NULL = unchecked |
| `status_code` | INT | HTTP status code | NULL = unchecked |
| `error_message` | TEXT | Error details | NULL |
| `created_at` | TIMESTAMP | Discovery time | DEFAULT CURRENT_TIMESTAMP |

**Indexes:**
- Primary: `id`
- Foreign key: `url_id` → `urls(id)` (CASCADE DELETE)
- Index: `is_internal` (for filtering)
- Index: `is_accessible` (for broken link queries)
- Index: `status_code` (for status filtering)

### 4. `api_tokens` - Authentication tokens

| Column | Type | Description | Constraints |
|--------|------|-------------|-------------|
| `id` | BIGINT | Primary key | AUTO_INCREMENT |
| `token_hash` | VARCHAR(255) | SHA256 hash of token | NOT NULL, UNIQUE |
| `name` | VARCHAR(100) | Token description | NOT NULL |
| `is_active` | BOOLEAN | Token status | DEFAULT TRUE |
| `expires_at` | TIMESTAMP | Expiration time | NULLABLE |
| `created_at` | TIMESTAMP | Creation time | DEFAULT CURRENT_TIMESTAMP |
| `last_used_at` | TIMESTAMP | Last usage time | NULLABLE |

**Indexes:**
- Primary: `id`
- Unique: `token_hash` (for fast lookup)
- Index: `is_active` (for filtering)
- Index: `expires_at` (for expiration checks)

## GORM Models

### URL Status Enum

```go
type URLStatus string

const (
    StatusQueued    URLStatus = "queued"    // Waiting to be crawled
    StatusRunning   URLStatus = "running"   // Currently being crawled
    StatusCompleted URLStatus = "completed" // Successfully crawled
    StatusError     URLStatus = "error"     // Crawl failed
)
```

### Model Relationships

```go
// One-to-One: URL → CrawlResult
type URL struct {
    CrawlResult *CrawlResult `gorm:"foreignKey:URLID"`
}

// One-to-Many: URL → FoundLinks
type URL struct {
    FoundLinks []FoundLink `gorm:"foreignKey:URLID"`
}
```

### Model Hooks

**BeforeCreate Hooks:**
- `URL.BeforeCreate()`: Sets default status to "queued"
- `CrawlResult.BeforeCreate()`: Sets crawled_at timestamp
- `APIToken.UpdateLastUsed()`: Updates last_used_at

### Business Logic Methods

**URL Model:**
- `TableName()`: Returns "urls"

**CrawlResult Model:**
- `GetHeadingCounts()`: Returns map of heading counts
- `GetTotalLinks()`: Returns sum of internal + external links

**FoundLink Model:**
- `IsBroken()`: Returns true if status code >= 400
- `GetStatusCategory()`: Returns human-readable status category

**APIToken Model:**
- `IsExpired()`: Checks if token is past expiration
- `IsValid()`: Checks if token is active and not expired
- `HashToken(token)`: Creates SHA256 hash

## Database Operations

### Connection Configuration

```go
// Environment Variables
DB_HOST=database
DB_PORT=3306  
DB_USER=crawler_user
DB_PASSWORD=crawler_password
DB_NAME=crawler_db
```

### Connection Features

- **Retry Logic**: 30 attempts with 2-second intervals
- **Connection Pooling**: 10 idle, 100 max connections
- **Health Checks**: Ping verification
- **Auto-Migration**: GORM handles schema updates

### Query Patterns

**Pagination:**
```go
query.Offset(offset).Limit(pageSize).Find(&urls)
```

**Filtering:**
```go
query.Where("status = ?", status)
query.Where("url LIKE ?", "%search%")
```

**Sorting:**
```go
query.Order("created_at DESC")
```

## Performance Considerations

### Indexing Strategy

1. **Primary Keys**: All tables have BIGINT AUTO_INCREMENT
2. **Foreign Keys**: Indexed for JOIN performance
3. **Search Fields**: url, status, created_at indexed
4. **Unique Constraints**: Prevent duplicate URLs and tokens

### Query Optimization

- Use pagination to limit result sets
- Index frequently filtered columns (status, is_internal)
- Preload relationships to avoid N+1 queries
- Use connection pooling for concurrent requests

## Migration Strategy

### Initial Setup

The database is initialized using:
1. Docker initialization script (`init.sql`)
2. GORM auto-migration for schema updates
3. Default data seeding (development token)

### Schema Updates

For production updates:
1. Create migration scripts
2. Test on staging environment  
3. Apply with zero-downtime strategies
4. Backup before major changes

## Backup and Recovery

### Development
- Docker volumes persist data between restarts
- Use `docker-compose down -v` to reset

### Production
- Regular MySQL backups using mysqldump
- Point-in-time recovery with binary logs
- Replica servers for high availability

## Security

### Access Control
- Dedicated database user with limited privileges
- No root access from application
- Connection over internal Docker network

### Data Protection
- Token hashes instead of plain text
- Prepared statements prevent SQL injection
- Input validation at application layer

## Monitoring

### Health Checks
- Database connectivity monitoring
- Connection pool status
- Query performance metrics
- Disk space monitoring

### Logging
- GORM query logging in development
- Silent mode in production
- Error logging for debugging