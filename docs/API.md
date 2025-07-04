# API Documentation

## Overview

The Web Crawler API is a RESTful service built with Go and Gin framework. It provides endpoints for URL management, authentication, and crawl result retrieval.

**Base URL:** `http://localhost:8080`  
**API Version:** v1  
**Authentication:** Bearer Token  

## Authentication

### Token-Based Authentication

All protected endpoints require a Bearer token in the Authorization header:

```http
Authorization: Bearer <token>
```

### Development Token

For development and testing:
```
Token: dev-token-12345
```

### Token Validation

**POST** `/api/auth/validate`

Validates an API token without requiring authentication.

**Request Body:**
```json
{
  "token": "dev-token-12345"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "valid": true,
    "token_name": "Development Token",
    "expires_at": null
  }
}
```

**Response (Invalid Token):**
```json
{
  "success": true,
  "data": {
    "valid": false
  }
}
```

### Current Token Info

**GET** `/api/auth/me`

Returns information about the current authenticated token.

**Headers:**
```http
Authorization: Bearer dev-token-12345
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "valid": true,
    "token_name": "Development Token",
    "expires_at": null
  }
}
```

## URL Management

### List URLs

**GET** `/api/urls`

Retrieves a paginated list of URLs with optional filtering and sorting.

**Headers:**
```http
Authorization: Bearer dev-token-12345
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number (min: 1) |
| `page_size` | integer | 20 | Items per page (min: 1, max: 100) |
| `search` | string | - | Search in URL field |
| `status` | string | - | Filter by status (queued, running, completed, error) |
| `sort_by` | string | created_at | Sort field (id, url, status, created_at, updated_at) |
| `sort_dir` | string | desc | Sort direction (asc, desc) |

**Example Request:**
```http
GET /api/urls?page=1&page_size=10&status=queued&sort_by=created_at&sort_dir=desc
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "url": "https://example.com",
      "status": "queued",
      "error_message": null,
      "created_at": "2025-07-04T13:00:00Z",
      "updated_at": "2025-07-04T13:00:00Z",
      "crawl_result": null
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

### Add URL

**POST** `/api/urls`

Adds a new URL for crawling.

**Headers:**
```http
Authorization: Bearer dev-token-12345
Content-Type: application/json
```

**Request Body:**
```json
{
  "url": "https://example.com"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "url": "https://example.com",
    "status": "queued",
    "error_message": null,
    "created_at": "2025-07-04T13:00:00Z",
    "updated_at": "2025-07-04T13:00:00Z"
  }
}
```

**Error Response (409 Conflict - Duplicate URL):**
```json
{
  "success": false,
  "error": {
    "code": "URL_EXISTS",
    "message": "URL already exists in the system"
  }
}
```

**Error Response (400 Bad Request - Invalid URL):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_URL",
    "message": "Invalid URL format",
    "details": "URL must use http or https protocol"
  }
}
```

### Get URL

**GET** `/api/urls/{id}`

Retrieves details of a specific URL.

**Headers:**
```http
Authorization: Bearer dev-token-12345
```

**Path Parameters:**
- `id` (integer): URL ID

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "url": "https://example.com",
    "status": "completed",
    "error_message": null,
    "created_at": "2025-07-04T13:00:00Z",
    "updated_at": "2025-07-04T13:05:00Z",
    "crawl_result": {
      "id": 1,
      "html_version": "HTML5",
      "page_title": "Example Domain",
      "heading_counts": {
        "h1": 1,
        "h2": 0,
        "h3": 0,
        "h4": 0,
        "h5": 0,
        "h6": 0
      },
      "internal_links_count": 3,
      "external_links_count": 2,
      "inaccessible_links_count": 0,
      "has_login_form": false,
      "crawled_at": "2025-07-04T13:05:00Z",
      "crawl_duration_ms": 1250,
      "total_links": 5
    }
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "URL_NOT_FOUND",
    "message": "URL not found"
  }
}
```

### Get URL Details

**GET** `/api/urls/{id}/details`

Retrieves comprehensive URL information including found links.

**Headers:**
```http
Authorization: Bearer dev-token-12345
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "url": "https://example.com",
    "status": "completed",
    "crawl_result": {
      "id": 1,
      "html_version": "HTML5",
      "page_title": "Example Domain",
      "heading_counts": {
        "h1": 1,
        "h2": 0,
        "h3": 0,
        "h4": 0,
        "h5": 0,
        "h6": 0
      },
      "internal_links_count": 3,
      "external_links_count": 2,
      "inaccessible_links_count": 1,
      "has_login_form": false,
      "total_links": 5
    },
    "found_links": [
      {
        "id": 1,
        "link_url": "https://example.com/about",
        "link_text": "About Us",
        "is_internal": true,
        "is_accessible": true,
        "status_code": 200,
        "error_message": null,
        "is_broken": false,
        "status_category": "success",
        "created_at": "2025-07-04T13:05:00Z"
      },
      {
        "id": 2,
        "link_url": "https://external-broken-link.com",
        "link_text": "Broken Link",
        "is_internal": false,
        "is_accessible": false,
        "status_code": 404,
        "error_message": "Not Found",
        "is_broken": true,
        "status_category": "client_error",
        "created_at": "2025-07-04T13:05:00Z"
      }
    ]
  }
}
```

### Delete URL

**DELETE** `/api/urls/{id}`

Deletes a URL and all associated data (cascading delete).

**Headers:**
```http
Authorization: Bearer dev-token-12345
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "URL deleted successfully",
    "id": 1
  }
}
```

### Bulk Delete URLs

**DELETE** `/api/urls/bulk`

Deletes multiple URLs in a single transaction.

**Headers:**
```http
Authorization: Bearer dev-token-12345
Content-Type: application/json
```

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "URLs deleted successfully",
    "deleted_count": 3,
    "ids": [1, 2, 3]
  }
}
```

## System Endpoints

### Health Check

**GET** `/health`

Returns system health status (no authentication required).

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2025-07-04T13:00:00Z",
  "service": "web-crawler-api"
}
```

### API Documentation

**GET** `/api`

Returns API documentation and endpoint overview (no authentication required).

**Response (200 OK):**
```json
{
  "service": "Web Crawler API",
  "version": "1.0.0",
  "endpoints": {
    "health": "GET /health",
    "auth": {
      "validate": "POST /api/auth/validate",
      "me": "GET /api/auth/me (auth required)"
    },
    "urls": {
      "list": "GET /api/urls (auth required)",
      "create": "POST /api/urls (auth required)",
      "get": "GET /api/urls/:id (auth required)",
      "details": "GET /api/urls/:id/details (auth required)",
      "delete": "DELETE /api/urls/:id (auth required)",
      "bulk_delete": "DELETE /api/urls/bulk (auth required)"
    }
  },
  "authentication": {
    "type": "Bearer Token",
    "header": "Authorization: Bearer <token>",
    "dev_token": "dev-token-12345"
  }
}
```

## Error Handling

### Standard Error Response

All errors follow the same response format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Additional error details (optional)"
  }
}
```

### HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, DELETE |
| 201 | Created | Successful POST |
| 400 | Bad Request | Invalid request format or validation error |
| 401 | Unauthorized | Missing or invalid authentication |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists (duplicate URL) |
| 500 | Internal Server Error | Server-side error |

### Common Error Codes

| Code | Description |
|------|-------------|
| `MISSING_AUTH_HEADER` | Authorization header not provided |
| `INVALID_AUTH_FORMAT` | Authorization header format incorrect |
| `INVALID_TOKEN` | Token is invalid or expired |
| `INVALID_REQUEST` | Request body format is invalid |
| `INVALID_PARAMS` | Query parameters are invalid |
| `INVALID_URL` | URL format validation failed |
| `URL_EXISTS` | URL already exists in system |
| `URL_NOT_FOUND` | URL ID not found |
| `DATABASE_ERROR` | Database operation failed |

## Rate Limiting

Currently no rate limiting is implemented. For production deployment, consider:

- Token-based rate limiting (requests per minute)
- IP-based rate limiting for public endpoints
- Different limits for different endpoint types

## CORS Configuration

The API supports cross-origin requests from:
- `http://localhost:3000` (development frontend)
- `http://localhost:80` (production frontend)

Allowed methods: GET, POST, PUT, DELETE, OPTIONS  
Allowed headers: Origin, Content-Type, Authorization

## Examples

### Complete Workflow Example

```bash
# 1. Validate token
curl -X POST http://localhost:8080/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token": "dev-token-12345"}'

# 2. Add a URL
curl -X POST http://localhost:8080/api/urls \
  -H "Authorization: Bearer dev-token-12345" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

# 3. List URLs
curl http://localhost:8080/api/urls \
  -H "Authorization: Bearer dev-token-12345"

# 4. Get URL details (after crawling)
curl http://localhost:8080/api/urls/1/details \
  -H "Authorization: Bearer dev-token-12345"

# 5. Delete URL
curl -X DELETE http://localhost:8080/api/urls/1 \
  -H "Authorization: Bearer dev-token-12345"
```

### Search and Pagination Example

```bash
# Search for URLs containing "example"
curl "http://localhost:8080/api/urls?search=example&page=1&page_size=5" \
  -H "Authorization: Bearer dev-token-12345"

# Filter by status and sort by creation date
curl "http://localhost:8080/api/urls?status=completed&sort_by=created_at&sort_dir=asc" \
  -H "Authorization: Bearer dev-token-12345"
```