package health

import "context"

type Service interface {
	GetHealthStatus(ctx context.Context) (*HealthCheckResponse, error)
}
