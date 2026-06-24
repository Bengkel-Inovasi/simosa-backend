package database

import (
	"fmt"
	"log"
	"os"
	"website/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbType := os.Getenv("DB_TYPE")

	if dbType == "postgres" {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				os.Getenv("DB_HOST"),
				os.Getenv("DB_USER"),
				os.Getenv("DB_PASSWORD"),
				os.Getenv("DB_NAME"),
				os.Getenv("DB_PORT"),
			)
		}
		fmt.Println("Connecting to PostgreSQL database...")
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		// Default to SQLite for local development
		fmt.Println("Using SQLite database (simosa.db)")
		DB, err = gorm.Open(sqlite.Open("simosa.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection established")

	// Auto Migrate the schemas
	err = DB.AutoMigrate(&models.Node{}, &models.SensorReading{}, &models.Harvest{}, &models.Expense{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed")
}
