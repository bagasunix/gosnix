package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
)

func SetupCustomerRoutes(router fiber.Router, handler *handlers.CustomerHandler) {
	router.Get("", handler.GetAllCustomer)
	router.Post("", handler.Create)
}
