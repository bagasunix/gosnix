package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

// Middleware generator
func LoggingMiddleware(logger *log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// jalankan handler berikutnya
		err := c.Next()

		// ambil response body
		resBody := string(c.Response().Body())

		// log ke console
		logger.Info().
			Str("method", c.Method()).
			Str("endpoint", c.OriginalURL()).
			Int("status", c.Response().StatusCode()).
			Str("user_agent", c.Get("User-Agent")).
			Str("ip_address", c.IP()).
			Dur("duration", time.Since(start)).
			Msg(resBody)

		return err
	}
}
