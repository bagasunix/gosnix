package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/entities"
	"github.com/bagasunix/gosnix/internal/domain/service"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/requests"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
	"github.com/bagasunix/gosnix/internal/infrastructure/persistence/postgres"
	"github.com/bagasunix/gosnix/internal/infrastructure/persistence/redis_client"
	"github.com/bagasunix/gosnix/pkg/configs"
	"github.com/bagasunix/gosnix/pkg/errors"
	"github.com/bagasunix/gosnix/pkg/hash"
	"github.com/bagasunix/gosnix/pkg/utils"
)

type customerService struct {
	cfg    *configs.Cfg
	repo   postgres.Repositories
	cache  redis_client.RedisClient
	logger *log.Logger
	// redis  *redis.Client
}

// Create implements customer.CustomerUsecase.
func (c *customerService) Create(ctx context.Context, request *requests.CreateCustomer) (response responses.BaseResponse[responses.CustomerResponse]) {
	if request.Validate() != nil {
		response.Code = 400
		response.Message = "phone atau password salah, silahkan coba lagi"
		response.Errors = request.Validate().Error()
		return response
	}
	phone, err := utils.ValidatePhone(request.Phone)
	if err != nil {
		response.Code = 409
		response.Message = "Validasi phone invalid"
		response.Errors = err.Error()
		return response
	}

	// Check if phone & email already exists
	checkData := c.repo.GetCustomer().FindByPhoneOrEmail(ctx, *phone, request.Email)
	if checkData.Value.Phone == *phone {
		response.Code = 409
		response.Message = "Phone dan email sudah terdaftar"
		response.Errors = "phone " + errors.ERR_ALREADY_EXISTS + ", email " + errors.ERR_ALREADY_EXISTS
		return response
	}
	if checkData.Error != nil && !strings.Contains(checkData.Error.Error(), "not found") {
		response.Code = 409
		response.Message = "Validasi phone invalid"
		response.Errors = checkData.Error.Error()
		return response
	}

	// Build customer
	customerBuild := new(entities.Customer)
	customerBuild.Name = request.Name
	customerBuild.Sex = int8(request.Sex)
	// parse DOB
	if request.DOB != nil {
		// dob, err := time.Parse("2006-01-02", *request.DOB)
		// if err != nil {
		// 	response.Code = 400
		// 	response.Message = "Tanggal lahir tidak valid, gunakan format YYYY-MM-DD"
		// 	response.Errors = err.Error()
		// 	return response
		// }
		// customerBuild.DOB = &dob
		customerBuild.DOB = request.DOB
	}
	customerBuild.Email = request.Email
	customerBuild.Phone = *phone
	customerBuild.PasswordHash = hash.HashAndSalt([]byte(request.Password)) // Pastikan untuk menghash password
	customerBuild.Address = request.Address
	customerBuild.Photo = request.Photo
	customerBuild.IsActive = 1

	tx := c.repo.GetCustomer().GetConnection().(*gorm.DB).Begin()

	if err = c.repo.GetCustomer().SaveTx(ctx, tx, customerBuild); err != nil {
		tx.Rollback()
		response.Code = 409
		response.Message = "Gagal membuat pelanggan"
		response.Errors = err.Error()
		return response
	}

	vehicles := make([]entities.Vehicle, 0, len(request.Vehicle))
	devices := make([]entities.DeviceGPS, 0, len(request.Vehicle))
	vehicleDevices := make([]entities.VehicleDevice, 0, len(request.Vehicle))

	if len(request.Vehicle) != 0 {
		platSet := make(map[string]struct{}, len(request.Vehicle))
		imeiSet := make(map[string]struct{}, len(request.Vehicle))

		for _, v := range request.Vehicle {
			checkVehicleCategory := c.repo.GetVehicleCategory().FindByParam(ctx, map[string]any{"id": v.CategoryID})
			if checkVehicleCategory.Error != nil {
				tx.Rollback()
				response.Code = 409
				response.Message = "Kategori kendaraan tidak ditemukan"
				response.Errors = checkVehicleCategory.Error.Error()
				return response
			}
			// Cek duplikat plat dalam request
			if _, exists := platSet[v.PlateNo]; exists {
				tx.Rollback()
				response.Code = 400
				response.Message = "Nomor plat duplikat dalam request"
				response.Errors = "Plat " + v.PlateNo + " sudah ada di request"
				return response
			}
			platSet[v.PlateNo] = struct{}{}

			// Cek duplikat IMEI dalam request
			if v.DeviceGPS.IMEI != "" {
				if _, exists := imeiSet[v.DeviceGPS.IMEI]; exists {
					tx.Rollback()
					response.Code = 400
					response.Message = "IMEI duplikat dalam request"
					response.Errors = "IMEI " + v.DeviceGPS.IMEI + " sudah ada di request"
					return response
				}
				imeiSet[v.DeviceGPS.IMEI] = struct{}{}

				// Cek IMEI di database
				checkImei := c.repo.GetDeviceGPS().FindByParam(ctx, map[string]interface{}{"imei": v.DeviceGPS.IMEI})
				if checkImei.Value.IMEI == v.DeviceGPS.IMEI {
					tx.Rollback()
					response.Code = 409
					response.Message = "Device GPS sudah terdaftar"
					response.Errors = "IMEI " + errors.ERR_ALREADY_EXISTS
					return response
				}
				if checkImei.Error != nil && !strings.Contains(checkImei.Error.Error(), "not found") {
					tx.Rollback()
					response.Code = 409
					response.Message = "Validasi IMEI invalid"
					response.Errors = checkImei.Error.Error()
					return response
				}
			}

			// Cek plat di database
			checkPlat := c.repo.GetVehicle().FindByParam(ctx, map[string]interface{}{"plate_no": v.PlateNo})
			if checkPlat.Value.PlateNo == v.PlateNo {
				tx.Rollback()
				response.Code = 409
				response.Message = "Kendaraan sudah terdaftar"
				response.Errors = "vehicle " + errors.ERR_ALREADY_EXISTS
				return response
			}
			if checkPlat.Error != nil && !strings.Contains(checkPlat.Error.Error(), "not found") {
				tx.Rollback()
				response.Code = 409
				response.Message = "Validasi vehicle invalid"
				response.Errors = checkPlat.Error.Error()
				return response
			}

			// Build vehicle
			vehicle := entities.Vehicle{
				Brand:           v.Brand,
				Color:           v.Color,
				CategoryID:      v.CategoryID,
				CustomerID:      customerBuild.ID,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				VIN:             v.VIN,
				ManufactureYear: v.ManufactureYear,
				CreatedBy:       customerBuild.ID,
			}
			vehicles = append(vehicles, vehicle)

			// Jika ada device GPS, buat device dan vehicle device
			if v.DeviceGPS.IMEI != "" {
				device := entities.DeviceGPS{
					IMEI:      v.DeviceGPS.IMEI,
					Brand:     v.DeviceGPS.Brand,
					Model:     v.DeviceGPS.Model,
					Protocol:  v.DeviceGPS.Protocol,
					SecretKey: v.DeviceGPS.SecretKey,
					CreatedBy: customerBuild.ID,
				}
				devices = append(devices, device)
			}
		}
	}

	// Simpan kendaraan
	if len(vehicles) != 0 {
		if err := c.repo.GetVehicle().SaveBatchTx(ctx, tx, vehicles); err != nil {
			tx.Rollback()
			return responses.BaseResponse[responses.CustomerResponse]{
				Code:    409,
				Message: "Gagal membuat kendaraan",
				Errors:  err.Error(),
			}
		}
	}

	// Simpan devices
	if len(devices) != 0 {
		if err := c.repo.GetDeviceGPS().SaveBatchTx(ctx, tx, devices); err != nil {
			tx.Rollback()
			return responses.BaseResponse[responses.CustomerResponse]{
				Code:    409,
				Message: "Gagal membuat device GPS",
				Errors:  err.Error(),
			}
		}

		// Buat hubungan vehicle-device
		for i, device := range devices {
			vehicleDevice := entities.VehicleDevice{
				VehicleID: vehicles[i].ID,
				DeviceID:  device.ID,
				StartTime: time.Now(),
				IsActive:  1,
			}
			vehicleDevices = append(vehicleDevices, vehicleDevice)
		}

		if err := c.repo.GetVehicleDevice().SaveBatchTx(ctx, tx, vehicleDevices); err != nil {
			tx.Rollback()
			return responses.BaseResponse[responses.CustomerResponse]{
				Code:    409,
				Message: "Gagal menghubungkan kendaraan dengan device",
				Errors:  err.Error(),
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		response.Code = 409
		response.Message = "Gagal membuat pelanggan"
		response.Errors = err.Error()
		return response
	}

	// Hapus cache Redis
	if err = c.cache.GetCustomerCache().DeleteByPattern(ctx, "*"); err != nil {
		response.Code = 409
		response.Message = "Gagal membuat pelanggan"
		response.Errors = err.Error()
		return response
	}

	// Build response
	resBuild := &responses.CustomerResponse{
		ID:        customerBuild.ID,
		Name:      customerBuild.Name,
		Sex:       customerBuild.Sex,
		Phone:     customerBuild.Phone,
		Email:     customerBuild.Email,
		Address:   customerBuild.Address,
		Photo:     customerBuild.Photo,
		IsActive:  customerBuild.IsActive,
		CreatedAt: customerBuild.CreatedAt,
	}

	if len(vehicles) > 0 {
		resBuild.Vehicle = make([]responses.VehicleResponse, 0, len(vehicles))
		for i, v := range vehicles {
			vehicleResponse := responses.VehicleResponse{
				ID:              v.ID,
				Brand:           v.Brand,
				Color:           v.Color,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
			}

			// Tambahkan informasi device jika ada
			if i < len(devices) {
				vehicleResponse.Device = &responses.DeviceGPSResponse{
					ID:          devices[i].ID,
					IMEI:        devices[i].IMEI,
					Brand:       devices[i].Brand,
					Model:       devices[i].Model,
					Protocol:    devices[i].Protocol,
					IsActive:    1,
					InstalledAt: time.Now(),
					CreatedAt:   devices[i].CreatedAt,
				}
			}

			resBuild.Vehicle = append(resBuild.Vehicle, vehicleResponse)
		}
	}

	response.Code = 201
	response.Message = "Sukses mendaftar"
	response.Data = resBuild
	return response
}

// DeleteCustomer implements customer.CustomerUsecase.
func (c *customerService) DeleteCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[any]) {
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Validasi error"
		response.Errors = request.Validate().Error()
		return response
	}

	intId, _ := strconv.Atoi(request.Id.(string))
	result := c.repo.GetCustomer().FindByParam(ctx, map[string]any{"id": intId})
	if result.Error != nil {
		response.Code = 400
		response.Message = "Data tidak ditemukan"
		response.Errors = result.Error.Error()
		return response
	}

	if err := c.repo.GetCustomer().Delete(ctx, intId); err != nil {
		response.Code = 400
		response.Message = "Data gagal dihapus"
		response.Errors = err.Error()
		return response
	}

	// 3. Hapus cache Redis
	// Hapus cache detail customer
	if err := c.cache.GetCustomerCache().Delete(ctx, intId); err != nil {
		response.Code = 400
		response.Message = "Data gagal dihapus"
		response.Errors = err.Error()
		return response
	}

	// Hapus cache list (opsional, pattern search)
	if err := c.cache.GetCustomerCache().DeleteByPattern(ctx, "search=*"); err != nil {
		response.Code = 400
		response.Message = "Data gagal dihapus"
		response.Errors = err.Error()
		return response
	}

	response.Message = "Data berhasil dihapus"
	response.Code = 200
	return response
}

// ListCustomer implements customer.CustomerUsecase.
func (c *customerService) ListCustomer(ctx context.Context, request *requests.BaseRequest) (response responses.BaseResponse[[]responses.CustomerResponse]) {
	// --- Validasi request
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Error Validasi"
		response.Errors = err.Error()
		return response
	}

	intPage, _ := strconv.Atoi(request.Page)
	intLimit, _ := strconv.Atoi(request.Limit)
	offset, limit := utils.CalculateOffsetAndLimit(intPage, intLimit)

	// --- Redis key
	cacheKey := fmt.Sprintf("data:search=%s:page=%d:limit=%d", request.Search, intPage, intLimit)
	countKey := fmt.Sprintf("count:search=%s", request.Search)

	var custResponse []responses.CustomerResponse

	// --- Cek Redis
	if val, err := c.cache.GetCustomerCache().Get(ctx, cacheKey); err == nil {
		if err := json.Unmarshal([]byte(val), &custResponse); err == nil {
			// Ambil totalItems dari cache count
			totalItems, _ := c.cache.GetCustomerCache().GetCount(ctx, countKey)
			totalPages := (totalItems + limit - 1) / limit
			response.Data = &custResponse
			response.Paging = &responses.PageMetadata{
				Page:      intPage,
				Size:      limit,
				TotalItem: totalItems,
				TotalPage: totalPages,
			}
			response.Message = "Inquiry pelanggan berhasil (cache)"
			response.Code = 200
			return response
		}
	}

	// --- Ambil dari DB
	resCust := c.repo.GetCustomer().FindAll(ctx, limit, offset, request.Search)
	if resCust.Error != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = resCust.Error.Error()
		return response
	}

	// --- Hitung total items
	totalItems, err := c.repo.GetCustomer().CountCustomer(ctx, request.Search)
	if err != nil {
		response.Code = 400
		response.Message = "Gagal menghitung data"
		response.Errors = err.Error()
		return response
	}
	totalPages := (totalItems + limit - 1) / limit

	// --- Mapping ke response
	custResponse = make([]responses.CustomerResponse, 0, len(resCust.Value))
	for _, v := range resCust.Value {
		vehResponse := make([]responses.VehicleResponse, 0, len(v.Vehicles)) // ✅ reset tiap customer

		for _, veh := range v.Vehicles {
			vResp := responses.VehicleResponse{
				ID:              veh.ID,
				Brand:           veh.Brand,
				Color:           veh.Color,
				Category:        veh.Category.Name,
				FuelType:        veh.FuelType,
				MaxSpeed:        veh.MaxSpeed,
				Model:           veh.Model,
				PlateNo:         veh.PlateNo,
				ManufactureYear: veh.ManufactureYear,
				IsActive:        veh.IsActive,
			}

			// --- mapping device aktif
			for _, vd := range veh.Devices {
				if vd.IsActive == 1 {
					vResp.Device = &responses.DeviceGPSResponse{
						ID:            vd.Device.ID,
						IMEI:          vd.Device.IMEI,
						Brand:         vd.Device.Brand,
						Model:         vd.Device.Model,
						Protocol:      vd.Device.Protocol,
						IsActive:      vd.IsActive,
						InstalledAt:   vd.StartTime,
						UninstalledAt: vd.EndTime,
						CreatedAt:     vd.Device.CreatedAt,
						UpdatedAt:     vd.UpdatedAt,
					}
					break
				}
			}

			vehResponse = append(vehResponse, vResp)
		}

		custResponse = append(custResponse, responses.CustomerResponse{
			ID:       v.ID,
			Name:     v.Name,
			Sex:      v.Sex,
			Email:    v.Email,
			Phone:    v.Phone,
			Address:  v.Address,
			Photo:    v.Photo,
			IsActive: v.IsActive,
			Vehicle:  vehResponse,
		})
	}

	// --- Simpan ke Redis (fallback kalau gagal)
	if data, err := json.Marshal(custResponse); err == nil {
		if err := c.cache.GetCustomerCache().Set(ctx, 5*time.Minute, data, cacheKey); err != nil {
			log.Printf("[WARN] gagal set cache data: %v", err)
		}
	}
	if err := c.cache.GetCustomerCache().Set(ctx, 5*time.Minute, totalItems, countKey); err != nil {
		log.Printf("[WARN] gagal set cache count: %v", err)
	}

	// --- Response akhir
	response.Data = &custResponse
	response.Paging = &responses.PageMetadata{
		Page:      intPage,
		Size:      limit,
		TotalItem: totalItems,
		TotalPage: totalPages,
	}
	response.Message = "Inquiry pelanggan berhasil"
	response.Code = 200
	return response
}

// UpdateCustomer implements customer.CustomerUsecase.
func (c *customerService) UpdateCustomer(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Validasi error"
		response.Errors = request.Validate().Error()
		return response
	}

	checkCust := c.repo.GetCustomer().FindByParam(ctx, map[string]any{"id": request.ID})
	if checkCust.Error != nil {
		response.Code = 409
		response.Message = "Pelanggan tidak ditemukan"
		response.Errors = checkCust.Error.Error()
		return response
	}

	mCustt := new(entities.Customer)
	// parse DOB
	if request.DOB != nil {
		// dob, err := time.Parse("2006-01-02", *request.DOB)
		// if err != nil {
		// 	response.Code = 400
		// 	response.Message = "Tanggal lahir tidak valid, gunakan format YYYY-MM-DD"
		// 	response.Errors = err.Error()
		// 	return response
		// }
		// mCustt.DOB = &dob
		mCustt.DOB = request.DOB
	}

	if request.Name != nil {
		mCustt.Name = *request.Name
	}
	if request.Sex != nil {
		mCustt.Sex = *request.Sex
	}
	if request.Photo != nil {
		mCustt.Photo = *request.Photo
	}
	if request.Address != nil {
		mCustt.Address = *request.Address
	}

	// --- Redis key berdasarkan customer ID
	intCustID, _ := strconv.Atoi(request.ID)

	if err := c.repo.GetCustomer().Updates(ctx, intCustID, mCustt); err != nil {
		response.Code = 400
		response.Message = "Gagal memperbarui pelanggan"
		response.Errors = err.Error()
		return response
	}

	resCust := new(responses.CustomerResponse)
	resCust.ID = checkCust.Value.ID
	resCust.Name = mCustt.Name
	resCust.Sex = mCustt.Sex
	resCust.DOB = mCustt.DOB
	resCust.Photo = mCustt.Photo
	resCust.Email = checkCust.Value.Email
	resCust.Phone = checkCust.Value.Phone
	resCust.Address = mCustt.Address
	resCust.IsActive = checkCust.Value.IsActive

	if len(checkCust.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCust.Value.Vehicles))
		for _, v := range checkCust.Value.Vehicles {
			resVec := responses.VehicleResponse{
				ID:              v.ID,
				Brand:           v.Brand,
				Color:           v.Color,
				Category:        v.Category.Name,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
			}

			for _, vd := range v.Devices {
				if vd.IsActive == 1 {
					resVec.Device = &responses.DeviceGPSResponse{
						ID:            vd.Device.ID,
						IMEI:          vd.Device.IMEI,
						Brand:         vd.Device.Brand,
						Model:         vd.Device.Model,
						Protocol:      vd.Device.Protocol,
						IsActive:      vd.IsActive,
						InstalledAt:   vd.StartTime,
						UninstalledAt: vd.EndTime,
						CreatedAt:     vd.Device.CreatedAt,
						UpdatedAt:     vd.UpdatedAt,
					}
					break
				}
			}
			resCust.Vehicle = append(resCust.Vehicle, resVec)
		}
	}

	// Hapus cache detail
	if err := c.cache.GetCustomerCache().Delete(ctx, intCustID); err != nil {
		response.Code = 400
		response.Message = "Gagal memperbarui pelanggan"
		response.Errors = err.Error()
		return response
	}

	// 4. Hapus cache list (opsional)
	if err := c.cache.GetCustomerCache().DeleteByPattern(ctx, "search=*"); err != nil {
		response.Code = 400
		response.Message = "Gagal memperbarui pelanggan"
		response.Errors = err.Error()
		return response
	}

	response.Data = &resCust
	response.Code = 200
	response.Message = "Berhsail memperbarui pelanggan"
	return response
}

// ViewCustomer implements customer.CustomerUsecase.
func (c *customerService) ViewCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.CustomerResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Validasi error"
		response.Errors = err.Error()
		return response
	}

	paramID, _ := strconv.Atoi(request.Id.(string))
	// --- Redis key berdasarkan customer ID

	//  Cek Redis dulu
	resCust := new(responses.CustomerResponse)
	if val, err := c.cache.GetCustomerCache().Get(ctx, paramID); err == nil {
		if err := json.Unmarshal([]byte(val), &resCust); err == nil {
			response.Data = &resCust
			response.Message = "Pelanggan ditemukan (cache)"
			response.Code = 200
			return response
		}
	}

	//  Ambil dari DB
	checkCustomer := c.repo.GetCustomer().FindByParam(ctx, map[string]any{"id": paramID})
	if checkCustomer.Error != nil {
		response.Code = 404
		response.Message = "Pelanggan tidak ditemukan"
		response.Errors = checkCustomer.Error.Error()
		return response
	}

	resCust.ID = checkCustomer.Value.ID
	resCust.Name = checkCustomer.Value.Name
	resCust.Email = checkCustomer.Value.Email
	resCust.Phone = checkCustomer.Value.Phone
	resCust.Address = checkCustomer.Value.Address
	resCust.Photo = checkCustomer.Value.Photo
	resCust.IsActive = checkCustomer.Value.IsActive

	if len(checkCustomer.Value.Vehicles) != 0 {
		resCust.Vehicle = make([]responses.VehicleResponse, 0, len(checkCustomer.Value.Vehicles))
		for _, v := range checkCustomer.Value.Vehicles {
			respV := responses.VehicleResponse{
				ID:              v.ID,
				Brand:           v.Brand,
				Color:           v.Color,
				Category:        v.Category.Name,
				FuelType:        v.FuelType,
				MaxSpeed:        v.MaxSpeed,
				Model:           v.Model,
				PlateNo:         v.PlateNo,
				ManufactureYear: v.ManufactureYear,
				IsActive:        v.IsActive,
			}

			for _, vd := range v.Devices {
				if vd.IsActive == 1 {
					respV.Device = &responses.DeviceGPSResponse{
						ID:            vd.Device.ID,
						IMEI:          vd.Device.IMEI,
						Brand:         vd.Device.Brand,
						Model:         vd.Device.Model,
						Protocol:      vd.Device.Protocol,
						IsActive:      vd.IsActive,
						InstalledAt:   vd.StartTime,
						UninstalledAt: vd.EndTime,
						CreatedAt:     vd.Device.CreatedAt,
						UpdatedAt:     vd.UpdatedAt,
					}
				}
			}

			resCust.Vehicle = append(resCust.Vehicle, respV)
		}
	}

	// Simpan ke Redis dengan expire 5 menit
	data, _ := json.Marshal(resCust)
	if err := c.cache.GetCustomerCache().Set(ctx, 5*time.Minute, data, paramID); err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}

	response.Data = &resCust
	response.Message = "Pelanggan ditemukan"
	response.Code = 200

	return response
}

// ViewCustomerWithVehicle implements customer.CustomerUsecase.
func (c *customerService) ViewCustomerWithVehicle(ctx context.Context, request *requests.BaseVehicle) (response responses.BaseResponse[[]responses.VehicleResponse]) {
	if err := request.Validate(); err != nil {
		response.Code = 400
		response.Message = "Validasi error"
		response.Errors = err.Error()
		return response
	}

	intPage, _ := strconv.Atoi(request.Page)
	intLimit, _ := strconv.Atoi(request.Limit)
	paramID, _ := strconv.Atoi(request.CustomerID)
	// --- Redis key berdasarkan customer ID
	cacheKey := fmt.Sprintf("customer:%v:search=%s:page=%d:limit=%d", paramID, request.Search, intPage, intLimit)
	countKey := fmt.Sprintf("customer:%s:count:search=%s", request.CustomerID, request.Search)
	offset, limit := utils.CalculateOffsetAndLimit(intPage, intLimit)

	//  Cek Redis dulu
	var resVehicle []responses.VehicleResponse
	// val, err := c.redis.Get(ctx, cacheKey).Result()
	val, err := c.cache.GetVehicleCache().Get(ctx, cacheKey)
	if err == nil {
		// Cache hit, unmarshal JSON
		if err := json.Unmarshal([]byte(*val), &resVehicle); err == nil {
			response.Data = &resVehicle
			response.Message = "Inquiry kendaraan berhasil"
			response.Code = 200
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
				ID:              v.ID,
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
	if err = c.cache.GetVehicleCache().Set(ctx, 5*time.Minute, data, cacheKey); err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}
	if err = c.cache.GetVehicleCache().Set(ctx, 5*time.Minute, totalItems, countKey); err != nil {
		response.Code = 400
		response.Message = "Gagal menarik data"
		response.Errors = err.Error()
		return response
	}

	response.Data = &resVehicle
	response.Paging = &responses.PageMetadata{
		Page:      intPage,
		Size:      limit,
		TotalItem: totalItems,
		TotalPage: totalPages,
	}
	response.Message = "Inquiry kendaraan berhasil"
	response.Code = 200
	return response
}

func NewCustomerService(logger *log.Logger, cache redis_client.RedisClient, repo postgres.Repositories, cfg *configs.Cfg) service.CustomerService {
	return &customerService{repo: repo, logger: logger, cfg: cfg, cache: cache}
}
