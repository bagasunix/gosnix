package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/entities"
	"github.com/bagasunix/gosnix/internal/domain/service"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/requests"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
	"github.com/bagasunix/gosnix/internal/infrastructure/persistence/postgres"
	"github.com/bagasunix/gosnix/pkg/configs"
	"github.com/bagasunix/gosnix/pkg/errors"
	"github.com/bagasunix/gosnix/pkg/utils"
)

type customerService struct {
	cfg    *configs.Cfg
	db     *gorm.DB
	repo   postgres.Repositories
	logger *log.Logger
	redis  *redis.Client
}

// Create implements customer.CustomerUsecase.
func (c *customerService) Create(ctx context.Context, request *requests.CreateCustomer) (response responses.BaseResponse[responses.CustomerResponse]) {
	if request.Validate() != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "phone atau password salah, silahkan coba lagi"
		response.Errors = request.Validate().Error()
		return response
	}
	phone, err := utils.ValidatePhone(request.Phone)
	if err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Validasi phone invalid"
		response.Errors = err.Error()
		return response
	}

	checkEmail := c.repo.GetCustomer().FindByParams(ctx, map[string]any{"phone": &phone})
	if checkEmail.Value.Phone == *phone {
		response.Code = fiber.StatusConflict
		response.Message = "Phone sudah terdaftar"
		response.Errors = "phone " + errors.ERR_ALREADY_EXISTS
		return response
	}
	if checkEmail.Error != nil && !strings.Contains(checkEmail.Error.Error(), "not found") {
		response.Code = fiber.StatusConflict
		response.Message = "Validasi phone invalid"
		response.Errors = checkEmail.Error.Error()
		return response
	}

	intSex, _ := strconv.Atoi(request.Sex)
	// Build customer
	customerBuild := &entities.Customer{
		Name:         request.Name,
		Sex:          int8(intSex),
		DOB:          request.DOB,
		Email:        request.Email,
		Phone:        *phone,
		PasswordHash: request.Password,
		Address:      request.Address,
		Photo:        request.Photo,
		IsActive:     1,
	}

	tx := c.repo.GetCustomer().GetConnection().(*gorm.DB).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = c.repo.GetCustomer().SaveTx(ctx, tx, customerBuild); err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pelanggan"
		response.Errors = err.Error()
		return response
	}

	vehicles := make([]entities.Vehicle, 0, len(request.Vehicle))
	if len(request.Vehicle) != 0 {
		platSet := make(map[string]struct{}, len(request.Vehicle))
		for _, v := range request.Vehicle {
			// Cek duplikat dalam request
			if _, exists := platSet[v.PlateNo]; exists {
				response.Code = fiber.StatusBadRequest
				response.Message = "Nomor plat duplikat dalam request"
				response.Errors = "Plat " + v.PlateNo + " sudah ada di request"
				return response
			}
			platSet[v.PlateNo] = struct{}{}

			checkPlat := c.repo.GetVehicle().FindByParam(ctx, map[string]interface{}{"plate_no": v.PlateNo})
			if checkPlat.Value.PlateNo == v.PlateNo {
				response.Code = fiber.StatusConflict
				response.Message = "Kendaraan sudah terdaftar"
				response.Errors = "vehicle " + errors.ERR_ALREADY_EXISTS
				return response
			}
			if checkPlat.Error != nil && !strings.Contains(checkPlat.Error.Error(), "not found") {
				response.Code = fiber.StatusConflict
				response.Message = "Validasi vehicle invalid"
				response.Errors = checkPlat.Error.Error()
				return response
			}

			vehicles = append(vehicles, entities.Vehicle{
				Brand:           v.Brand,
				Color:           v.Color,
				CustomerID:      customerBuild.ID,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				CreatedBy:       customerBuild.ID,
			})
		}
	}

	if len(vehicles) != 0 {
		if err := c.repo.GetVehicle().SaveBatchTx(ctx, tx, vehicles); err != nil {
			tx.Rollback()
			return responses.BaseResponse[responses.CustomerResponse]{
				Code:    fiber.StatusConflict,
				Message: "Gagal membuat pelanggan",
				Errors:  err.Error(),
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Gagal membuat pelanggan"
		response.Errors = err.Error()
		return response
	}

	// Hapus cache Redis
	// hapus semua key yang dimulai dengan "customers:"
	keys, _ := c.redis.Keys(ctx, "customers:*").Result()
	if len(keys) > 0 {
		c.redis.Del(ctx, keys...)
	}

	// Build response
	resBuild := &responses.CustomerResponse{
		ID:       strconv.Itoa(customerBuild.ID),
		Name:     customerBuild.Name,
		Phone:    customerBuild.Phone,
		Email:    customerBuild.Email,
		Address:  customerBuild.Address,
		IsActive: strconv.Itoa(int(customerBuild.IsActive)),
	}

	if len(vehicles) > 0 {
		resBuild.Vehicle = make([]responses.VehicleResponse, 0, len(vehicles))
		for _, v := range vehicles {
			resBuild.Vehicle = append(resBuild.Vehicle, responses.VehicleResponse{
				ID:              v.ID.String(),
				Brand:           v.Brand,
				Color:           v.Color,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
			})
		}
	}

	response.Code = fiber.StatusOK
	response.Message = "Sukses mendaftar"
	response.Data = resBuild
	return response
}

// DeleteCustomer implements customer.CustomerUsecase.
func (c *customerService) DeleteCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[any]) {
	if err := request.Validate(); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = request.Validate().Error()
		return response
	}

	intId, _ := strconv.Atoi(request.Id.(string))
	result := c.repo.GetCustomer().FindByParams(ctx, map[string]any{"id": intId})
	if result.Error != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Data tidak ditemukan"
		response.Errors = result.Error.Error()
		return response
	}

	if err := c.repo.GetCustomer().Delete(ctx, intId); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Data gagal dihapus"
		response.Errors = err.Error()
		return response
	}

	// 3. Hapus cache Redis
	cacheKey := fmt.Sprintf("customers:%d", intId)
	c.redis.Del(ctx, cacheKey) // Hapus cache detail

	// Hapus cache list (opsional, pattern search)
	listKeys, _ := c.redis.Keys(ctx, "customers:search=*").Result()
	if len(listKeys) > 0 {
		c.redis.Del(ctx, listKeys...)
	}

	response.Message = "Data berhasil dihapus"
	response.Code = fiber.StatusOK
	return response
}

// ListCustomer implements customer.CustomerUsecase.
func (c *customerService) ListCustomer(ctx context.Context, request *requests.BaseRequest) (response responses.BaseResponse[[]responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Error Validasi"
		response.Errors = err.Error()
		return response
	}
	intPage, _ := strconv.Atoi(request.Page)
	intLimit, _ := strconv.Atoi(request.Limit)

	offset, limit := utils.CalculateOffsetAndLimit(intPage, intLimit)

	// --- Redis key berdasarkan search, page, limit
	cacheKey := fmt.Sprintf("customers:search=%s:page=%d:limit=%d", request.Search, intPage, intLimit)

	var custResponse []responses.CustomerResponse

	// Cek Redis
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal langsung
		if err := json.Unmarshal([]byte(val), &custResponse); err == nil {
			// Hit, kembalikan response dengan paging
			totalItems, _ := c.redis.Get(ctx, "customers:count:search="+request.Search).Int() // optional cache count
			totalPages := (totalItems + limit - 1) / limit
			response.Data = &custResponse
			response.Paging = &responses.PageMetadata{
				Page:      intPage,
				Size:      limit,
				TotalItem: totalItems,
				TotalPage: totalPages,
			}
			response.Message = "Inquiry pelanggan berhasil"
			response.Code = fiber.StatusOK
			return response
		}
	}

	//  Ambil dari DB
	resCust := c.repo.GetCustomer().FindAll(ctx, limit, offset, request.Search)
	if resCust.Error != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = resCust.Error.Error()
		return response
	}
	// Calculate total items and total pages
	totalItems, err := c.repo.GetCustomer().CountCustomer(ctx, request.Search)
	if err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}
	totalPages := (totalItems + limit - 1) / limit

	// Map ke response
	custResponse = make([]responses.CustomerResponse, 0, len(resCust.Value))
	for _, v := range resCust.Value {
		custResponse = append(custResponse, responses.CustomerResponse{
			ID:       strconv.Itoa(v.ID),
			Name:     v.Name,
			Email:    v.Email,
			Phone:    v.Phone,
			Address:  v.Address,
			IsActive: strconv.Itoa(int(v.IsActive)),
		})
	}

	// Simpan ke Redis
	data, _ := json.Marshal(custResponse)
	c.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	c.redis.Set(ctx, "customers:count:search="+request.Search, totalItems, 5*time.Minute)

	response.Data = &custResponse
	response.Paging = &responses.PageMetadata{
		Page:      intPage,
		Size:      limit,
		TotalItem: totalItems,
		TotalPage: totalPages,
	}
	response.Message = "Inquiry pelanggan berhasil"
	response.Code = fiber.StatusOK
	return response
}

// UpdateCustomer implements customer.CustomerUsecase.
func (c *customerService) UpdateCustomer(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = request.Validate().Error()
		return response
	}

	checkCust := c.repo.GetCustomer().FindByParams(ctx, map[string]any{"id": request.ID})
	if checkCust.Error != nil {
		response.Code = fiber.StatusConflict
		response.Message = "Pelanggan tidak ditemukan"
		response.Errors = checkCust.Error.Error()
		return response
	}

	mCustt := new(entities.Customer)
	mCustt.Name = request.Name
	mCustt.Phone = request.Phone
	mCustt.Address = request.Address

	// --- Redis key berdasarkan customer ID
	intCustID, _ := strconv.Atoi(request.ID)
	cacheKey := fmt.Sprintf("customers:%d", intCustID)

	if err := c.repo.GetCustomer().Updates(ctx, intCustID, mCustt); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Gagal memperbarui pelanggan"
		response.Errors = err.Error()
		return response
	}

	resCust := new(responses.CustomerResponse)
	resCust.ID = strconv.Itoa(mCustt.ID)
	resCust.Name = mCustt.Name
	resCust.Email = mCustt.Email
	resCust.Phone = mCustt.Phone
	resCust.Address = mCustt.Address
	resCust.IsActive = strconv.Itoa(int(mCustt.IsActive))

	if len(checkCust.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCust.Value.Vehicles))
		for _, v := range checkCust.Value.Vehicles {
			resCust.Vehicle = append(resCust.Vehicle, responses.VehicleResponse{
				Brand:           v.Brand,
				Color:           v.Color,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
			})
		}
	}

	// Hapus cache detail
	c.redis.Del(ctx, cacheKey)

	// 4. Hapus cache list (opsional)
	listKeys, _ := c.redis.Keys(ctx, "customers:search=*").Result()
	if len(listKeys) > 0 {
		c.redis.Del(ctx, listKeys...)
	}

	response.Data = &resCust
	response.Code = fiber.StatusOK
	response.Message = "Berhsail memperbarui pelanggan"
	return response
}

// ViewCustomer implements customer.CustomerUsecase.
func (c *customerService) ViewCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.CustomerResponse]) {
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

	checkCustomer := c.repo.GetCustomer().FindByParams(ctx, map[string]any{"id": paramID})
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
	resCust.IsActive = strconv.Itoa(int(checkCustomer.Value.IsActive))

	if len(checkCustomer.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCustomer.Value.Vehicles))
		for _, v := range checkCustomer.Value.Vehicles {
			resCust.Vehicle = append(resCust.Vehicle, responses.VehicleResponse{
				Brand:           v.Brand,
				Color:           v.Color,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
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
func (c *customerService) ViewCustomerWithVehicle(ctx context.Context, request *requests.BaseVehicle) (response responses.BaseResponse[[]responses.VehicleResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = fiber.StatusBadRequest
		response.Message = "Validasi error"
		response.Errors = err.Error()
		return response
	}

	intPage, _ := strconv.Atoi(request.Page)
	intLimit, _ := strconv.Atoi(request.Limit)
	paramID, _ := strconv.Atoi(request.CustomerID)
	// --- Redis key berdasarkan customer ID
	cacheKey := fmt.Sprintf("vehicles:customer:%s:search=%s:page=%d:limit=%d", paramID, request.Search, intPage, intLimit)
	countKey := fmt.Sprintf("vehicles:customer:%s:count:search=%s", request.CustomerID, request.Search)
	offset, limit := utils.CalculateOffsetAndLimit(intPage, intLimit)

	//  Cek Redis dulu
	var resVehicle []responses.VehicleResponse
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit, unmarshal JSON
		if err := json.Unmarshal([]byte(val), &resVehicle); err == nil {
			response.Data = &resVehicle
			response.Message = "Inquiry kendaraan berhasil"
			response.Code = fiber.StatusOK
			return response
		}
	}

	//  Ambil dari DB
	checkVehicle := c.repo.GetVehicle().FindByCustomer(ctx, paramID, limit, offset, request.Search)
	if checkVehicle.Error != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = checkVehicle.Error.Error()
		return response
	}
	// Calculate total items and total pages
	totalItems, err := c.repo.GetVehicle().CountVehicle(ctx, request.Search)
	if err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}
	totalPages := (totalItems + limit - 1) / limit

	if len(checkVehicle.Value) != 0 {
		resVehicle = make([]responses.VehicleResponse, 0, len(checkVehicle.Value))
		for _, v := range checkVehicle.Value {
			resVehicle = append(resVehicle, responses.VehicleResponse{
				ID:              v.ID.String(),
				Brand:           v.Brand,
				Color:           v.Color,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
			})
		}
	}

	// Simpan ke Redis dengan expire 5 menit
	data, _ := json.Marshal(resVehicle)
	c.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	c.redis.Set(ctx, countKey, totalItems, 5*time.Minute)

	response.Data = &resVehicle
	response.Paging = &responses.PageMetadata{
		Page:      intPage,
		Size:      limit,
		TotalItem: totalItems,
		TotalPage: totalPages,
	}
	response.Message = "Inquiry kendaraan berhasil"
	response.Code = fiber.StatusOK
	return response
}

func NewCustomerService(logger *log.Logger, db *gorm.DB, repo postgres.Repositories, redis *redis.Client, cfg *configs.Cfg) service.CustomerService {
	return &customerService{db: db, repo: repo, logger: logger, redis: redis, cfg: cfg}
}
