package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"avoids-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// PostgreSQL connection parameters
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Create connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		host, port, user, password, dbname)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Optional: add additional configuration
		// TranslateError: true, // Translate database errors to more meaningful messages
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Avoid{}, &models.DailyCheck{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Optional: Add database connection pool configuration
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
