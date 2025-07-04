# Backend Documentation

## Overview

The backend is built with Go 1.22 using the Gin web framework, GORM for database operations, and follows clean architecture principles with clear separation of concerns.

## Architecture

### Project Structure

```
backend/
├── main.go                 # Application entry point
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
│   └── urls.go          # URL management handlers
├── middleware/
│   └── auth.go          # Authentication middleware
└── integration_test.go   # Integration tests
```

### Design Patterns

**Repository Pattern**: Models encapsulate database operations  
**DTO Pattern**: Separate request/response objects from domain models  
**Middleware Pattern**: Cross-cutting concerns like authentication  
**Handler Pattern**: HTTP request processing and response formatting

## Core Components

### Main Application (`main.go`)

The application entry point that:
- Initializes database connection
- Configures Gin router and middleware
- Sets up CORS for frontend integration
- Defines route groups and handlers
- Starts the HTTP server

```go
func main() {
    // Database initialization
    database.Connect()
    
    // Router setup
    router := gin.New()
    router.Use(gin.Logger(), gin.Recovery())
    
    // CORS configuration
    router.Use(cors.New(corsConfig))
    
    // Route groups
    public