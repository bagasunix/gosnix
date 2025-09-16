package service

import (
	"context"

	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
)

type HealthService interface {
	GetHealthStatus(ctx context.Context) (*responses.HealthCheckResponse, error)
}
