package handlers

import (
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

// Create godoc
// @Summary Membuat customer baru
// @Description Membuat data customer baru beserta kendaraan opsional
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   request body requests.CreateCustomer true "Customer Request"
// @Success 201 {object} responses.CustomerRegisterResponseWrapper
// @Failure 400 {object} responses.ErrorBadRequestResponse
// @Failure 401 {object} responses.ErrorUnauthorizedResponse
// @Failure 409 {object} responses.ErrorConflictResponse
// @Router /customers [post]
// @Security BearerAuth
func (ac *CustomerHandler) Create(ctx *fiber.Ctx) error {
	request := new(requests.CreateCustomer)
	var response responses.BaseResponse[responses.CustomerResponse]
	if err := ctx.BodyParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response = ac.service.Create(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

// GetAllCustomer godoc
// @Summary Mendapatkan daftar customer
// @Description Menampilkan daftar customer dengan pagination dan optional search
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   page   query string false "Page number"
// @Param   limit  query string false "Limit per page"
// @Param   search query string false "Search by name/email"
// @Success 200 {object} responses.CustomerListResponseWrapper
// @Failure 400 {object} responses.ErrorBadRequestResponse
// @Failure 401 {object} responses.ErrorUnauthorizedResponse
// @Router /customers [get]
// @Security BearerAuth
func (c *CustomerHandler) GetAllCustomer(ctx *fiber.Ctx) error {
	request := new(requests.BaseRequest)
	var response responses.BaseResponse[[]responses.CustomerResponse]

	if err := ctx.QueryParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	response = c.service.ListCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}
	return ctx.Status(response.Code).JSON(response)
}
