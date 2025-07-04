package handlers

import (
	"net/http"

	"web-crawler/database"
	"web-crawler/dto"
	"web-crawler/models"
	
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct{}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// ValidateToken validates an API token
// POST /api/auth/validate
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req dto.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"INVALID_REQUEST",
			"Invalid request format",
			err.Error(),
		))
		return
	}
	
	// Validate the token using the same logic as middleware
	tokenHash := models.HashToken(req.Token)
	
	var apiToken models.APIToken
	result := database.DB.Where("token_hash = ?", tokenHash).First(&apiToken)
	
	if result.Error != nil {
		c.JSON(http.StatusOK, dto.SuccessResponse(dto.TokenValidationResponse{
			Valid: false,
		}))
		return
	}
	
	// Check if token is valid
	if !apiToken.IsValid() {
		c.JSON(http.StatusOK, dto.SuccessResponse(dto.TokenValidationResponse{
			Valid: false,
		}))
		return
	}
	
	// Token is valid
	c.JSON(http.StatusOK, dto.SuccessResponse(dto.TokenValidationResponse{
		Valid:     true,
		TokenName: apiToken.Name,
		ExpiresAt: apiToken.ExpiresAt,
	}))
}

// GetCurrentToken returns information about the current authenticated token
// GET /api/auth/me
func (h *AuthHandler) GetCurrentToken(c *gin.Context) {
	// Get token from context (set by auth middleware)
	tokenInterface, exists := c.Get("api_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(
			"NO_TOKEN",
			"No authenticated token found",
			"",
		))
		return
	}
	
	apiToken, ok := tokenInterface.(*models.APIToken)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"INTERNAL_ERROR",
			"Failed to process token information",
			"",
		))
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse(dto.TokenValidationResponse{
		Valid:     true,
		TokenName: apiToken.Name,
		ExpiresAt: apiToken.ExpiresAt,
	}))
}