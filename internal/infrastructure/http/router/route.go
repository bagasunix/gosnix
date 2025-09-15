package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
)

// SetupRoutes sekarang hanya menerima handler yang sudah di-inject
func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthHandler) {
	api := app.Group("/api")
	SetupHealthRoutes(api, healthHandler)
}
