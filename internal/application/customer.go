package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/repository"
	"github.com/bagasunix/gosnix/internal/domain/service"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/requests"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
)

type customerUsecase struct {
	db     *gorm.DB
	repo   repository.CustomerRepository
	logger *log.Logger
	redis  *redis.Client
}

// Create implements customer.CustomerUsecase.
func (c *customerUsecase) Create(ctx context.Context, request *requests.CreateCustomer) (response responses.BaseResponse[responses.CustomerResponse]) {
	panic("unimplemented")
}

// DeleteCustomer implements customer.CustomerUsecase.
func (c *customerUsecase) DeleteCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[any]) {
	panic("unimplemented")
}

// ListCustomer implements customer.CustomerUsecase.
func (c *customerUsecase) ListCustomer(ctx context.Context, request *requests.BaseRequest) (response responses.BaseResponse[[]responses.CustomerResponse]) {
	panic("unimplemented")
}

// UpdateCustomer implements customer.CustomerUsecase.
func (c *customerUsecase) UpdateCustomer(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.CustomerResponse]) {
	panic("unimplemented")
}

// ViewCustomer implements customer.CustomerUsecase.
func (c *customerUsecase) ViewCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = err.Error()
		return response
	}

	paramID, _ := strconv.Atoi(request.Id.(string))
	// --- Redis key berdasarkan customer ID
	cacheKey := fmt.Sprintf("customers:%d", paramID)

	//  Cek Redis dulu
	resCust := new(responses.CustomerResponse)
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal JSON
		if err := json.Unmarshal([]byte(val), &resCust); err == nil {
			response.Data = &resCust
			response.Message = "Pelanggan ditemukan"
			response.Code = fiber.StatusOK
			return response
		}
	}

	checkCustomer := c.repo.FindByParams(ctx, map[string]any{"id": paramID})
	if checkCustomer.Error != nil {
		response.Code = fiber.StatusNotFound
		response.Message = "Pelanggan tidak ditemukan"
		response.Errors = checkCustomer.Error.Error()
		return response
	}

	resCust.ID = strconv.Itoa(checkCustomer.Value.ID)
	resCust.Name = checkCustomer.Value.Name
	resCust.Email = checkCustomer.Value.Email
	resCust.Phone = checkCustomer.Value.Phone
	resCust.Address = checkCustomer.Value.Address
	resCust.IsActive = checkCustomer.Value.IsActive

	if len(checkCustomer.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCustomer.Value.Vehicles))
		for _, v := range checkCustomer.Value.Vehicles {
			resCust.Vehicle = append(resCust.Vehicle, responses.VehicleResponse{
				Brand:    v.Brand,
				Color:    v.Color,
				FuelType: v.FuelType,
				MaxSpeed: v.MaxSpeed,
				Model:    v.Model,
				PlateNo:  v.PlateNo,
				Year:     v.Year,
				IsActive: v.IsActive,
			})
		}
	}

	// Simpan ke Redis dengan expire 5 menit
	data, _ := json.Marshal(resCust)
	c.redis.Set(ctx, cacheKey, data, 5*time.Minute)

	response.Data = &resCust
	response.Message = "Pelanggan ditemukan"
	response.Code = fiber.StatusOK

	return response
}

// ViewCustomerWithVehicle implements customer.CustomerUsecase.
func (c *customerUsecase) ViewCustomerWithVehicle(ctx context.Context, request *requests.BaseVehicle) (response responses.BaseResponse[[]responses.VehicleResponse]) {
	panic("unimplemented")
}

func NewCustomerUsecase(logger *log.Logger, db *gorm.DB, repo repository.CustomerRepository, redis *redis.Client) service.CustomerUsecase {
	return &customerUsecase{db: db, repo: repo, logger: logger, redis: redis}
}
