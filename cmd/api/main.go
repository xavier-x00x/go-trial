package main

import (
	"log"

	"go-trial/internal/config"
	"go-trial/internal/delivery/http/route"
	"go-trial/internal/infrastructure/database"
	"go-trial/internal/registry"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db := database.NewMySQL(&cfg.Database)

	// Connect to Redis
	rdb := database.NewRedis(&cfg.Redis)

	// Initialize dependency injection
	reg := registry.New(db, rdb, cfg)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Go Trial API",
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Static files
	app.Static("/uploads", "./uploads")

	// Register routes
	route.Setup(app, reg, reg.JWTManager)

	// Start server
	addr := ":" + cfg.App.Port
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": err.Error(),
	})
}
