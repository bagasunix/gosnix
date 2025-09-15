package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	// _ "github.com/bagasunix/gosnix/docs"
	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
)

func SetupHealthRoutes(router fiber.Router, handler *handlers.HealthHandler) {
	router.Get("health", handler.CheckHealth)
	router.Get("swagger/*", swagger.HandlerDefault)
}
