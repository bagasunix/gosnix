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

// ViewCustomer godoc
// @Summary Mendapatkan detail customer
// @Description Menampilkan detail customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Success 200 {object} responses.CustomerDetailResponseWrapper
// @Failure 400 {object} responses.ErrorBadRequestResponse
// @Failure 404 {object} responses.ErrorNotFoundResponse
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

// UpdateCustomer godoc
// @Summary Melakukan update customer
// @Description Melakukan update customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Param   body body requests.UpdateCustomer true "Update Customer Request"
// @Success 200 {object} responses.CustomerUpdateResponseWrapper
// @Failure 400 {object} responses.ErrorBadRequestResponse
// @Failure 404 {object} responses.ErrorNotFoundResponse
// @Router /customers/{id} [put]
// @Security BearerAuth
func (c *CustomerHandler) UpdateCustomer(ctx *fiber.Ctx) error {
	request := new(requests.UpdateCustomer)
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

	if err := ctx.BodyParser(request); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Invalid request"
		response.Errors = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}
	request.ID = id
	response = c.service.UpdateCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}

// DeleteCustomer godoc
// @Summary Melakukan delete customer
// @Description Melakukan delete customer berdasarkan ID
// @Tags Customer
// @Accept  json
// @Produce  json
// @Param   id   path int true "Customer ID"
// @Success 200 {object} responses.CustomerDeleteResponseWrapper
// @Failure 400 {object} responses.ErrorBadRequestResponse
// @Failure 404 {object} responses.ErrorNotFoundResponse
// @Router /customers/{id} [delete]
// @Security BearerAuth
func (c *CustomerHandler) DeleteCustomer(ctx *fiber.Ctx) error {
	request := new(requests.EntityId)
	var response responses.BaseResponse[any]

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
	response = c.service.DeleteCustomer(ctx.Context(), request)
	if response.Errors != "" {
		return ctx.Status(response.Code).JSON(response)
	}

	return ctx.Status(response.Code).JSON(response)
}
