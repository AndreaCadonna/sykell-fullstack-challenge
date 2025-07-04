# Development Documentation

## Overview

This guide covers the development workflow, tools, and best practices for contributing to the Web Crawler project.

## Getting Started

### Prerequisites

- **Docker & Docker Compose**: Container runtime and orchestration
- **Git**: Version control
- **Code Editor**: VS Code recommended with extensions
- **Node.js 20+**: For local frontend development (optional)
- **Go 1.22+**: For local backend development (optional)

### Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd web-crawler

# Start development environment
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Verify services are running
curl http://localhost:8080/health  # Backend health check
curl http://localhost:3000         # Frontend
```

### Development Environment

**Services Available**:
- Frontend: http://localhost:3000 (React + Vite + Hot Reload)
- Backend: http://localhost:8080 (Go + Air + Hot Reload)
- Database: localhost:3306 (MySQL accessible for debugging)

**Authentication**:
- Development token: `dev-token-12345`
- Use in Authorization header: `Bearer dev-token-12345`

## Development Workflow

### Branch Strategy

```bash
# Feature development
git checkout -b feature/crawl-engine
git checkout -b bugfix/auth-token-validation
git checkout -b docs/api-documentation

# Commit convention
git commit -m "feat: add URL crawling functionality"
git commit -m "fix: resolve token validation issue"
git commit -m "docs: update API documentation"
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
feat: new feature
fix: bug fix
docs: documentation changes
style: formatting changes
refactor: code refactoring
test: adding tests
chore: maintenance tasks
```

### Development Commands

**Start Development**:
```bash
# Full stack development
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Individual services
docker-compose up frontend  # Frontend only
docker-compose up backend   # Backend only
docker-compose up database  # Database only
```

**Stop and Clean**:
```bash
# Stop services
docker-compose down

# Remove volumes (reset database)
docker-compose down -v

# Clean images and cache
docker system prune -a
```

## Frontend Development

### Hot Reload Setup

The frontend uses Vite with hot module replacement:

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    host: "0.0.0.0",  // Docker compatibility
    port: 3000,
    watch: {
      usePolling: true,  // File watching in Docker
    },
  },
});
```

### Component Development

**Create New Component**:
```bash
# Create component file
touch src/components/URLList.tsx

# Create test file
touch src/components/URLList.test.tsx

# Create story file (future)
touch src/components/URLList.stories.tsx
```

**Component Template**:
```typescript
import React from 'react';

interface URLListProps {
  urls: URL[];
  onSelect: (url: URL) => void;
}

export function URLList({ urls, onSelect }: URLListProps) {
  return (
    <div className="space-y-4">
      {urls.map((url) => (
        <div
          key={url.id}
          className="p-4 border rounded-lg cursor-pointer hover:bg-gray-50"
          onClick={() => onSelect(url)}
        >
          <h3 className="font-semibold">{url.url}</h3>
          <p className="text-sm text-gray-600">Status: {url.status}</p>
        </div>
      ))}
    </div>
  );
}

export default URLList;
```

### Testing Frontend

```bash
# Run tests in watch mode
npm run test

# Run tests once (CI mode)
npm run test:run

# Run tests with coverage
npm run coverage

# Run tests in Docker
docker-compose -f docker-compose.yml -f docker-compose.test.yml up frontend
```

**Test Example**:
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, test, expect, vi } from 'vitest';
import URLList from './URLList';

describe('URLList', () => {
  const mockUrls = [
    { id: 1, url: 'https://example.com', status: 'queued' },
    { id: 2, url: 'https://test.com', status: 'completed' },
  ];

  test('renders URLs correctly', () => {
    const onSelect = vi.fn();
    render(<URLList urls={mockUrls} onSelect={onSelect} />);
    
    expect(screen.getByText('https://example.com')).toBeInTheDocument();
    expect(screen.getByText('Status: queued')).toBeInTheDocument();
  });

  test('calls onSelect when URL clicked', () => {
    const onSelect = vi.fn();
    render(<URLList urls={mockUrls} onSelect={onSelect} />);
    
    fireEvent.click(screen.getByText('https://example.com'));
    expect(onSelect).toHaveBeenCalledWith(mockUrls[0]);
  });
});
```

## Backend Development

### Hot Reload Setup

The backend uses Air for automatic rebuilds:

```toml
# .air.toml
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  delay = 1000
```

### Adding New Endpoints

**1. Define Model** (if needed):
```go
// models/new_entity.go
package models

type NewEntity struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
}

func (NewEntity) TableName() string {
    return "new_entities"
}
```

**2. Create DTOs**:
```go
// dto/requests.go
type CreateNewEntityRequest struct {
    Name string `json:"name" binding:"required"`
}

func (r *CreateNewEntityRequest) Validate() error {
    if len(r.Name) < 3 {
        return fmt.Errorf("name must be at least 3 characters")
    }
    return nil
}
```

**3. Implement Handler**:
```go
// handlers/new_entity.go
package handlers

func (h *NewEntityHandler) Create(c *gin.Context) {
    var req dto.CreateNewEntityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, dto.ErrorResponse("INVALID_REQUEST", "Invalid request format", err.Error()))
        return
    }

    if err := req.Validate(); err != nil {
        c.JSON(400, dto.ErrorResponse("VALIDATION_ERROR", "Validation failed", err.Error()))
        return
    }

    entity := models.NewEntity{Name: req.Name}
    if err := database.DB.Create(&entity).Error; err != nil {
        c.JSON(500, dto.ErrorResponse("DATABASE_ERROR", "Failed to create entity", err.Error()))
        return
    }

    c.JSON(201, dto.SuccessResponse(entity))
}
```

**4. Register Route**:
```go
// main.go
protected.POST("/entities", newEntityHandler.Create)
```

### Testing Backend

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestURLValidation ./models

# Run tests in Docker
docker-compose -f docker-compose.yml -f docker-compose.test.yml up backend
```

**Test Example**:
```go
// handlers/new_entity_test.go
func TestCreateNewEntity(t *testing.T) {
    // Setup test database
    setupTestDB(t)
    router := setupTestRouter()

    // Test data
    requestBody := dto.CreateNewEntityRequest{
        Name: "Test Entity",
    }
    jsonBody, _ := json.Marshal(requestBody)

    // Make request
    req, _ := http.NewRequest("POST", "/api/entities", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer test-token")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assertions
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response dto.APIResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.True(t, response.Success)
}
```

## Database Development

### Schema Changes

**1. Update Model**:
```go
// Add new field to existing model
type URL struct {
    // ... existing fields
    Priority int `json:"priority" gorm:"default:0"`
}
```

**2. Run Migration**:
```bash
# Auto-migration on restart
docker-compose restart backend

# Manual migration (future)
docker exec backend-container ./migrate up
```

**3. Update Tests**:
```go
func TestURLWithPriority(t *testing.T) {
    url := models.URL{
        URL: "https://example.com",
        Priority: 5,
    }
    // Test logic
}
```

### Database Debugging

```bash
# Connect to database
docker exec -it <database-container> mysql -u crawler_user -pcrawler_password crawler_db

# View tables
SHOW TABLES;

# View table structure
DESCRIBE urls;

# View data
SELECT * FROM urls LIMIT 10;

# View relationships
SELECT u.url, cr.page_title 
FROM urls u 
LEFT JOIN crawl_results cr ON u.id = cr.url_id;
```

### Database Seeding

**Development Data**:
```go
// database/seed.go (future)
func SeedDevelopmentData() {
    urls := []models.URL{
        {URL: "https://example.com", Status: models.StatusQueued},
        {URL: "https://test.com", Status: models.StatusCompleted},
    }
    
    for _, url := range urls {
        database.DB.Create(&url)
    }
}
```

## API Development

### Testing API Endpoints

**Manual Testing**:
```bash
# Health check
curl http://localhost:8080/health

# List URLs
curl -H "Authorization: Bearer dev-token-12345" \
     http://localhost:8080/api/urls

# Create URL
curl -X POST \
     -H "Authorization: Bearer dev-token-12345" \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com"}' \
     http://localhost:8080/api/urls
```

**Automated Testing Script**:
```bash
# Run comprehensive API tests
./test-api.sh

# Or use the provided test script
chmod +x test-api.sh
./test-api.sh
```

### API Documentation

**Update API docs** when adding endpoints:
```markdown
# docs/API.md

### New Endpoint

**POST** `/api/entities`

Creates a new entity.

**Request**:
```json
{
  "name": "Entity Name"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Entity Name",
    "created_at": "2025-07-04T13:00:00Z"
  }
}
```

## Code Quality

### Linting and Formatting

**Frontend**:
```bash
# ESLint
npm run lint
npm run lint -- --fix

# TypeScript check
npx tsc --noEmit
```

**Backend**:
```bash
# Go formatting
go fmt ./...

# Go linting (install golangci-lint)
golangci-lint run

# Go vet
go vet ./...
```

### Code Review Guidelines

**Before Creating PR**:
- [ ] All tests pass
- [ ] Code is formatted and linted
- [ ] Documentation updated
- [ ] No debug code or console.logs
- [ ] Environment variables documented
- [ ] Error handling implemented

**Review Checklist**:
- [ ] Code follows project conventions
- [ ] Security considerations addressed
- [ ] Performance implications considered
- [ ] Tests cover new functionality
- [ ] Documentation is clear and complete

## Debugging

### Frontend Debugging

**Browser DevTools**:
- React DevTools extension
- Network tab for API calls
- Console for JavaScript errors
- Performance tab for optimization

**Vite Debugging**:
```bash
# Debug mode
DEBUG=vite:* npm run dev

# Build analysis
npm run build -- --debug
```

### Backend Debugging

**Go Debugging**:
```go
// Add debug prints
import "log"
log.Printf("Debug: %+v", variable)

// Use delve debugger (future)
dlv debug
```

**Database Debugging**:
```go
// Enable GORM logging
db.Debug().Create(&model)

// Log SQL queries
logger := logger.Default.LogMode(logger.Info)
```

**Docker Debugging**:
```bash
# View container logs
docker-compose logs backend

# Enter container for debugging
docker exec -it backend-container sh

# Check environment variables
docker exec backend-container env
```

## Performance

### Frontend Performance

**Bundle Analysis**:
```bash
# Analyze bundle size
npm run build
npx vite-bundle-analyzer dist

# Check for unused code
npx unimported
```

**Performance Monitoring**:
```typescript
// Add performance timing
const start = performance.now();
// ... code to measure
const end = performance.now();
console.log(`Operation took ${end - start} milliseconds`);
```

### Backend Performance

**Profiling**:
```go
// Add pprof endpoint (development only)
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

**Database Performance**:
```sql
-- Analyze slow queries
SHOW PROCESSLIST;

-- Check query execution plan
EXPLAIN SELECT * FROM urls WHERE status = 'queued';

-- Monitor index usage
SHOW INDEX FROM urls;
```

## Environment Management

### Environment Variables

**Development** (`.env.development`):
```bash
DB_HOST=database
DB_PORT=3306
DB_USER=crawler_user
DB_PASSWORD=crawler_password
DB_NAME=crawler_db
JWT_SECRET=dev-secret-key
ENV=development
```

**Testing** (`.env.test`):
```bash
DB_NAME=crawler_test_db
JWT_SECRET=test-secret-key
ENV=testing
```

### Configuration Management

**Backend Config**:
```go
// config/config.go
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    JWTSecret  string
    Environment string
}

func LoadConfig() Config {
    return Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        // ... other fields
    }
}
```

## Deployment

### Local Deployment Testing

```bash
# Test production build locally
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build

# Test with production-like data
./scripts/seed-production-data.sh
```

### Pre-deployment Checklist

- [ ] All tests pass in CI environment
- [ ] Database migrations tested
- [ ] Environment variables configured
- [ ] SSL certificates ready (production)
- [ ] Backup procedures tested
- [ ] Monitoring setup verified
- [ ] Load testing completed (if needed)

## Troubleshooting

### Common Development Issues

**Docker Issues**:
```bash
# Containers won't start
docker-compose logs <service>
docker system df  # Check disk space

# Permission errors
sudo chown -R $USER:$USER .
chmod -R 755 .

# Port conflicts
lsof -i :3000  # Check what's using port
```

**Hot Reload Issues**:
```bash
# Frontend not reloading
export CHOKIDAR_USEPOLLING=true
export WATCHPACK_POLLING=true

# Backend not reloading
docker-compose restart backend
```

**