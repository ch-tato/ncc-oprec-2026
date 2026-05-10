package main

import (
	"log"

	"notes-app/internal/database"
	"notes-app/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialise the database (SQLite + GORM auto-migration).
	database.Init()

	// Create a Gin engine with default middleware (logger + recovery).
	r := gin.Default()

	// Load HTML templates.
	r.LoadHTMLGlob("templates/*")

	// Serve static assets (CSS, JS, images).
	r.Static("/static", "./static")

	// Register application routes.
	routes.Setup(r)

	// Start the server.
	log.Println("Server starting on http://localhost:8880")
	if err := r.Run(":8880"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
