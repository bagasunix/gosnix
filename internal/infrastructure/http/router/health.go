package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	// ini wajib, sesuaikan dengan nama module di go.mod kamu

	_ "github.com/bagasunix/gosnix/docs"
	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
)

func SetupHealthRoutes(router fiber.Router, handler *handlers.HealthHandler) {
	router.Get("health", handler.CheckHealth)
	// route swagger UI
	router.Use("/swagger", func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src * 'unsafe-inline' 'unsafe-eval'")
		return c.Next()
	})
	router.Get("swagger/*", swagger.HandlerDefault)
}
