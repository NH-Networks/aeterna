package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // <--- Nieuwe import
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// We kijken of er een DB_TYPE is ingesteld (standaard postgres)
	dbType := os.Getenv("DB_TYPE")
	dsn := os.Getenv("DATABASE_URL")
	var err error

	switch dbType {
	case "sqlite":
		// Als er geen bestandsnaam is opgegeven, gebruiken we 'aeterna.db'
		if dsn == "" {
			dsn = "aeterna.db"
		}
		// Verbinden met SQLite
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	default:
		// Standaard PostgreSQL logica (fallback)
		if dsn == "" {
			dsn = "host=localhost user=postgres password=postgres dbname=aeterna port=5432 sslmode=disable"
		}
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// BELANGRIJK: pgcrypto is specifiek voor Postgres.
	// We voeren dit alleen uit als we NIET op sqlite zitten.
	if dbType != "sqlite" {
		if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;").Error; err != nil {
			log.Fatal("Failed to enable pgcrypto extension: ", err)
		}
	}

	log.Println("Database connection successfully opened using driver:", dbType)
}
