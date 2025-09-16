package configs

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"

	"github.com/bagasunix/gosnix/internal/infrastructure/middlewares"
	"github.com/bagasunix/gosnix/pkg/configs"
)

func InitFiber(ctx context.Context, cfg *configs.Cfg, redis *redis.Client, logger *log.Logger) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		UnescapePath: true,
		ErrorHandler: middlewares.NewErrorHandler(logger),
		// EnablePrintRoutes: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		return c.Next()
	})
	app.Use(helmet.New())
	app.Use(recover.New())
	app.Use(favicon.New())
	app.Use(middlewares.Limitter(redis, cfg))
	app.Use(middlewares.LoggingMiddleware(logger))

	return app
}
