package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	// ini wajib, sesuaikan dengan nama module di go.mod kamu

	_ "github.com/bagasunix/gosnix/docs"
	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
	"github.com/bagasunix/gosnix/pkg/configs"
)

// SetupRoutes sekarang hanya menerima handler yang sudah di-inject
func SetupRoutes(app *fiber.App, cfg *configs.Cfg, handlerHealth *handlers.HealthHandler, handlerCustomer *handlers.CustomerHandler) {
	// route swagger UI
	app.Group(cfg.Server.Version).Use("/swagger", func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src * 'unsafe-inline' 'unsafe-eval'")
		return c.Next()
	})
	app.Group(cfg.Server.Version).Get("/swagger/*", swagger.HandlerDefault)
	SetupHealthRoutes(app.Group(cfg.Server.Version), handlerHealth)
	SetupCustomerRoutes(app.Group(cfg.Server.Version+"/customers"), handlerCustomer)
}
