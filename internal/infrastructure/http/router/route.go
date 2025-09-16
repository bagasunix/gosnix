package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	// ini wajib, sesuaikan dengan nama module di go.mod kamu

	_ "github.com/bagasunix/gosnix/docs"
	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
)

// SetupRoutes sekarang hanya menerima handler yang sudah di-inject
func SetupRoutes(app *fiber.App, handlerHealth *handlers.HealthHandler, handlerCustomer *handlers.CustomerHandler) {
	api := app.Group("/api")
	// route swagger UI
	api.Use("/swagger", func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src * 'unsafe-inline' 'unsafe-eval'")
		return c.Next()
	})
	api.Get("swagger/*", swagger.HandlerDefault)
	SetupHealthRoutes(api, handlerHealth)
	SetupCustomerRoutes(api, handlerCustomer)
}
