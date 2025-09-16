package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/bagasunix/gosnix/internal/domain/service"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/requests"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// ViewCustomer godoc
// @Summary Mendapatkan detail customer
// @Description Menampilkan detail customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Success 200 {object} responses.BaseResponseCustomer
// @Failure 400 {object} responses.BaseResponseSwagger
// @Failure 404 {object} responses.BaseResponseSwagger
// @Router /customers/{id} [get]
// @Security BearerAuth
func (c *CustomerHandler) ViewCustomer(ctx *fiber.Ctx) error {
	request := new(requests.EntityId)
	var response responses.BaseResponse[*responses.CustomerResponse]

	id := ctx.Params("id")
	if id == "" {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID tidak ditemukan"
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	if _, err := strconv.Atoi(id); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = "ID harus berupa angka"
		return ctx.Status(response.Code).JSON(response)
	}

	request.Id = id
	response = c.service.ViewCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}
