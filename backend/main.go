package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web-crawler/database"
	"web-crawler/handlers"
	"web-crawler/middleware"
	"web-crawler/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting Web Crawler API...")

	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize and start crawl manager
	crawlManager := services.NewCrawlManager()
	crawlManager.Start()
	defer crawlManager.Stop()

	// Set up graceful shutdown
	setupGracefulShutdown(crawlManager)

	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:80"}, // Frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler()
	urlHandler := handlers.NewURLHandler()
	crawlHandler := handlers.NewCrawlHandler(crawlManager)

	// Health check endpoint (no auth required)
	router.GET("/health", func(c *gin.Context) {
		response := gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "web-crawler-api",
		}

		// Add crawl manager status
		if crawlManager != nil {
			response["crawl_manager"] = crawlManager.GetQueueStatus()
		}

		c.JSON(http.StatusOK, response)
	})

	// Public routes (no authentication required)
	public := router.Group("/api")
	{
		public.POST("/auth/validate", authHandler.ValidateToken)
	}

	// Protected routes (authentication required)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth routes
		auth := protected.Group("/auth")
		{
			auth.GET("/me", authHandler.GetCurrentToken)
		}

		// URL routes
		urls := protected.Group("/urls")
		{
			urls.GET("", urlHandler.ListURLs)
			urls.POST("", urlHandler.AddURL)
			urls.GET("/:id", urlHandler.GetURL)
			urls.GET("/:id/details", urlHandler.GetURLDetails)
			urls.DELETE("/:id", urlHandler.DeleteURL)
			urls.DELETE("/bulk", urlHandler.BulkDeleteURLs)

			// Crawl control routes
			urls.POST("/:id/crawl", crawlHandler.StartCrawl)
			urls.GET("/:id/crawl/status", crawlHandler.GetCrawlStatus)
		}

		// Crawl management routes
		crawls := protected.Group("/crawls")
		{
			crawls.POST("/bulk", crawlHandler.StartBulkCrawl)
			crawls.GET("/queue/status", crawlHandler.GetQueueStatus)
		}
	}

	// API documentation endpoint
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Web Crawler API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health": "GET /health",
				"auth": gin.H{
					"validate": "POST /api/auth/validate",
					"me":       "GET /api/auth/me (auth required)",
				},
				"urls": gin.H{
					"list":         "GET /api/urls (auth required)",
					"create":       "POST /api/urls (auth required)",
					"get":          "GET /api/urls/:id (auth required)",
					"details":      "GET /api/urls/:id/details (auth required)",
					"delete":       "DELETE /api/urls/:id (auth required)",
					"bulk_delete":  "DELETE /api/urls/bulk (auth required)",
					"start_crawl":  "POST /api/urls/:id/crawl (auth required)",
					"crawl_status": "GET /api/urls/:id/crawl/status (auth required)",
				},
				"crawls": gin.H{
					"bulk_crawl":   "POST /api/crawls/bulk (auth required)",
					"queue_status": "GET /api/crawls/queue/status (auth required)",
				},
			},
			"authentication": gin.H{
				"type":      "Bearer Token",
				"header":    "Authorization: Bearer <token>",
				"dev_token": "dev-token-12345",
			},
			"crawl_manager": crawlManager.GetQueueStatus(),
		})
	})

	// Handle 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "endpoint not found",
			"code":  "NOT_FOUND",
		})
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Development token: dev-token-12345")
	log.Printf("API documentation: http://localhost:%s/api", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("Crawl manager status: %+v", crawlManager.GetQueueStatus())

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupGracefulShutdown configures graceful shutdown handling
func setupGracefulShutdown(crawlManager *services.CrawlManager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal...")

		// Stop crawl manager gracefully
		log.Println("Stopping crawl manager...")
		crawlManager.Stop()

		// Close database connection
		log.Println("Closing database connection...")
		database.Close()

		log.Println("Graceful shutdown completed")
		os.Exit(0)
	}()
}
