package application

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bagasunix/gosnix/internal/domain/health"
)

type healthService struct {
	repo health.Repository
}

func NewHealthService(repo health.Repository) health.Service {
	return &healthService{repo: repo}
}

func (s *healthService) GetHealthStatus(ctx context.Context) (*health.HealthCheckResponse, error) {
	response := &health.HealthCheckResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Timestamp: time.Now().Unix(),
		Details:   make(map[string]string),
	}

	var unhealthy bool

	// pakai errgroup buat paralel check
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.repo.CheckDB(ctx); err != nil {
			unhealthy = true
			response.Database = err.Error()
		} else {
			response.Database = "connected"
		}
		return nil
	})

	g.Go(func() error {
		if err := s.repo.CheckRedis(ctx); err != nil {
			unhealthy = true
			response.Redis = err.Error()
		} else {
			response.Redis = "connected"
		}
		return nil
	})

	g.Go(func() error {
		if err := s.repo.CheckRabbitMQ(ctx); err != nil {
			unhealthy = true
			response.RabbitMQ = err.Error()
		} else {
			response.RabbitMQ = "connected"
		}
		return nil
	})

	// tunggu semua goroutine selesai
	_ = g.Wait()

	if unhealthy {
		response.Status = "unhealthy"
	}

	return response, nil
}
