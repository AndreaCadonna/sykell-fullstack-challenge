package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"web-crawler/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// GetConfigFromEnv reads database configuration from environment variables
func GetConfigFromEnv() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "3306"),
		User:     getEnvWithDefault("DB_USER", "crawler_user"),
		Password: getEnvWithDefault("DB_PASSWORD", "crawler_password"),
		Name:     getEnvWithDefault("DB_NAME", "crawler_db"),
	}
}

// Connect establishes database connection and runs migrations
func Connect() error {
	config := GetConfigFromEnv()

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	// Configure GORM logger
	gormLogger := logger.Default
	if os.Getenv("ENV") == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Retry connection with backoff
	var err error
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})

		if err == nil {
			// Test connection
			sqlDB, err := DB.DB()
			if err == nil {
				if err := sqlDB.Ping(); err == nil {
					break // Connection successful
				}
			}
		}

		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")

	// Run auto-migration
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

// runMigrations runs GORM auto-migrations for all models
func runMigrations() error {
	log.Println("Running database migrations...")

	// Auto-migrate all models
	err := DB.AutoMigrate(
		&models.URL{},
		&models.CrawlResult{},
		&models.FoundLink{},
		&models.APIToken{},
	)

	if err != nil {
		return fmt.Errorf("auto-migration failed: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
