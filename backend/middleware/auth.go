package middleware

import (
	"net/http"
	"strings"

	"web-crawler/database"
	"web-crawler/models"
	
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates API tokens from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}
		
		// Parse Bearer token
		token := extractBearerToken(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format. Use: Bearer <token>",
				"code":  "INVALID_AUTH_FORMAT",
			})
			c.Abort()
			return
		}
		
		// Validate token
		apiToken, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			c.Abort()
			return
		}
		
		// Update last used timestamp
		apiToken.UpdateLastUsed()
		database.DB.Save(apiToken)
		
		// Store token in context for use in handlers
		c.Set("api_token", apiToken)
		c.Next()
	}
}

// extractBearerToken extracts token from "Bearer <token>" format
func extractBearerToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

// validateToken checks if the provided token is valid
func validateToken(token string) (*models.APIToken, error) {
	// Hash the token for database lookup
	tokenHash := models.HashToken(token)
	
	var apiToken models.APIToken
	result := database.DB.Where("token_hash = ?", tokenHash).First(&apiToken)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Check if token is valid (active and not expired)
	if !apiToken.IsValid() {
		return nil, result.Error
	}
	
	return &apiToken, nil
}

// OptionalAuthMiddleware is like AuthMiddleware but doesn't require authentication
// Useful for endpoints that can work with or without auth
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			token := extractBearerToken(authHeader)
			if token != "" {
				if apiToken, err := validateToken(token); err == nil {
					apiToken.UpdateLastUsed()
					database.DB.Save(apiToken)
					c.Set("api_token", apiToken)
				}
			}
		}
		c.Next()
	}
}