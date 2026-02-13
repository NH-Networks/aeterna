package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // <--- Added SQLite driver
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Check if DB_TYPE is set (default to postgres)
	dbType := os.Getenv("DB_TYPE")
	dsn := os.Getenv("DATABASE_URL")
	var err error

	switch dbType {
	case "sqlite":
		// If no filename is provided, use 'aeterna.db' as default
		if dsn == "" {
			dsn = "aeterna.db"
		}
		// Connect to SQLite
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	default:
		// Default PostgreSQL logic (fallback)
		if dsn == "" {
			dsn = "host=localhost user=postgres password=postgres dbname=aeterna port=5432 sslmode=disable"
		}
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// IMPORTANT: pgcrypto is specific to Postgres.
	// Only execute this if NOT using SQLite.
	if dbType != "sqlite" {
		// Ensure UUID generation is available for Postgres
		if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;").Error; err != nil {
			log.Fatal("Failed to enable pgcrypto extension: ", err)
		}
	}

	log.Println("Database connection successfully opened using driver:", dbType)
}
