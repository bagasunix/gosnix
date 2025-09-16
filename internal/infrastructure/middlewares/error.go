package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

func NewErrorHandler(logger *log.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		start := time.Now()
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		logger.Error().
			Str("method", c.Method()).
			Str("endpoint", c.OriginalURL()).
			Int("status", code).
			Str("ip_address", c.IP()).
			Err(err).
			Dur("duration", time.Since(start)).Msg(message)

		return c.Status(code).JSON(fiber.Map{
			"message": message,
			"errors":  err.Error(),
		})
	}
}
