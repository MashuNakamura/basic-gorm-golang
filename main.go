package main

import (
	"gorm-management-users/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize SQLite connection
	db, err := gorm.Open(sqlite.Open("db/data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate Database Schema
	migrateDatabase(db)

	// Initialize Fiber app
	app := fiber.New()

	// Register user routes
	routes.UserRoutes(app, db)

	// Simple ping endpoint
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Start server
	log.Println("Server listening on :3000")

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
