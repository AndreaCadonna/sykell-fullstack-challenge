package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"web-crawler/database"
	"web-crawler/dto"
	"web-crawler/handlers"
	"web-crawler/middleware"
	"web-crawler/models"
	
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB sets up an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	// Run migrations
	err = database.DB.AutoMigrate(
		&models.URL{},
		&models.CrawlResult{},
		&models.FoundLink{},
		&models.APIToken{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	
	// Create test API token
	testToken := models.APIToken{
		TokenHash: models.HashToken("test-token"),
		Name:      "Test Token",
		IsActive:  true,
	}
	database.DB.Create(&testToken)
}

// setupTestRouter creates a test Gin router with authentication
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add auth middleware to protected routes
	authHandler := handlers.NewAuthHandler()
	urlHandler := handlers.NewURLHandler()
	
	// Public routes
	public := router.Group("/api")
	{
		public.POST("/auth/validate", authHandler.ValidateToken)
	}
	
	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/auth/me", authHandler.GetCurrentToken)
		protected.GET("/urls", urlHandler.ListURLs)
		protected.POST("/urls", urlHandler.AddURL)
		protected.GET("/urls/:id", urlHandler.GetURL)
		protected.DELETE("/urls/:id", urlHandler.DeleteURL)
	}
	
	return router
}

func TestAuthValidateToken_HappyPath(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()
	
	// Test valid token
	requestBody := dto.ValidateTokenRequest{
		Token: "test-token",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	req, _ := http.NewRequest("POST", "/api/auth/validate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	var response dto.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if !response.Success {
		t.Error("Expected successful response")
	}
	
	// Check response data
	responseData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected response data to be a map")
		return
	}
	
	if valid, ok := responseData["valid"].(bool); !ok || !valid {
		t.Error("Expected token to be valid")
	}
}

func TestURLsCreate_HappyPath(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()
	
	// Test creating a URL
	requestBody := dto.AddURLRequest{
		URL: "https://example.com",
	}
	jsonBody, _ := json.Marshal(requestBody)
	
	req, _ := http.NewRequest("POST", "/api/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	
	var response dto.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if !response.Success {
		t.Error("Expected successful response")
	}
	
	// Verify URL was created in database
	var url models.URL
	result := database.DB.Where("url = ?", "https://example.com").First(&url)
	if result.Error != nil {
		t.Errorf("URL should be created in database: %v", result.Error)
	}
	
	if url.Status != models.StatusQueued {
		t.Errorf("Expected status to be %s, got %s", models.StatusQueued, url.Status)
	}
}

func TestURLsList_HappyPath(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()
	
	// Create test URLs
	testURLs := []models.URL{
		{URL: "https://example1.com", Status: models.StatusQueued},
		{URL: "https://example2.com", Status: models.StatusCompleted},
		{URL: "https://example3.com", Status: models.StatusError},
	}
	
	for _, url := range testURLs {
		database.DB.Create(&url)
	}
	
	// Test listing URLs
	req, _ := http.NewRequest("GET", "/api/urls?page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	var response dto.APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if !response.Success {
		t.Error("Expected successful response")
	}
	
	if response.Meta == nil {
		t.Error("Expected meta information for pagination")
	}
	
	if response.Meta.Total != 3 {
		t.Errorf("Expected total count to be 3, got %d", response.Meta.Total)
	}
}

func TestAuthenticationRequired(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()
	
	// Test accessing protected endpoint without token
	req, _ := http.NewRequest("GET", "/api/urls", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestInvalidToken(t *testing.T) {
	setupTestDB(t)
	router := setupTestRouter()
	
	// Test accessing protected endpoint with invalid token
	req, _ := http.NewRequest("GET", "/api/urls", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}