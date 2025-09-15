package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/gosnix/internal/domain/health"
)

type HealthHandler struct {
	service health.Service
}

func NewHealthHandler(service health.Service) *HealthHandler {
	return &HealthHandler{service: service}
}

func (h *HealthHandler) CheckHealth(c *fiber.Ctx) error {
	ctx := c.Context()
	status, err := h.service.GetHealthStatus(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if status.Status == "unhealthy" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(status)
	}

	return c.JSON(status)
}
