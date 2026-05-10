package database

import (
	"log"
	"os"

	"notes-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB holds the global database connection.
var DB *gorm.DB

// Init opens the SQLite connection, runs auto-migrations, and
// stores the resulting *gorm.DB in the package-level DB variable.
func Init() {
	// Ensure the data directory exists.
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open("data/notes.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the Note model.
	if err := DB.AutoMigrate(&models.Note{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database initialised successfully")
}
