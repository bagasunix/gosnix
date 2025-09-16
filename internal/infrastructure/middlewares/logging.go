package middlewares

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
)

func LoggingMiddleware(logger *log.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		// ambil response body
		resBody := string(c.Response().Body())
		status := c.Response().StatusCode()
		// jika body kosong, jangan log, biarkan ErrorHandler menangani
		if len(resBody) == 0 {
			return err
		}
		if strings.Contains(c.OriginalURL(), "swagger") {
			resBody = "Swagger Config"
		}

		// pilih level log berdasarkan status code
		entry := logger.Info()
		if status >= 400 && status < 500 {
			entry = logger.Warn()
		} else if status >= 500 {
			entry = logger.Error()
		}

		entry.
			Str("method", c.Method()).
			Str("endpoint", c.OriginalURL()).
			Int("status", status).
			Str("user_agent", c.Get("User-Agent")).
			Str("ip_address", c.IP()).
			Dur("duration", time.Since(start)).Msg(resBody)
		return err
	}
}
