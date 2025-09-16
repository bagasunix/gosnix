package service

import (
	"context"

	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/requests"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
)

type CustomerService interface {
	Create(ctx context.Context, request *requests.CreateCustomer) (response responses.BaseResponse[responses.CustomerResponse])
	ListCustomer(ctx context.Context, request *requests.BaseRequest) (response responses.BaseResponse[[]responses.CustomerResponse])
	ViewCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[*responses.CustomerResponse])
	DeleteCustomer(ctx context.Context, request *requests.EntityId) (response responses.BaseResponse[any])
	UpdateCustomer(ctx context.Context, request *requests.UpdateCustomer) (response responses.BaseResponse[*responses.CustomerResponse])
	ViewCustomerWithVehicle(ctx context.Context, request *requests.BaseVehicle) (response responses.BaseResponse[[]responses.VehicleResponse])
}
