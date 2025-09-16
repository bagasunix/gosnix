package application

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bagasunix/gosnix/internal/domain/service"
	"github.com/bagasunix/gosnix/internal/infrastructure/dtos/responses"
	"github.com/bagasunix/gosnix/internal/infrastructure/messaging/rabbitmq"
	"github.com/bagasunix/gosnix/internal/infrastructure/persistence/postgres"
	redis "github.com/bagasunix/gosnix/internal/infrastructure/persistence/redis_client"
)

type healthService struct {
	postgres    postgres.Repositories
	redisClient redis.RedisClient
	rmq         rabbitmq.RmqClient
}

func NewHealthService(postgres postgres.Repositories, redisClient redis.RedisClient, rmq rabbitmq.RmqClient) service.HealthService {
	return &healthService{postgres: postgres, redisClient: redisClient, rmq: rmq}
}

func (s *healthService) GetHealthStatus(ctx context.Context) (*responses.HealthCheckResponse, error) {
	response := &responses.HealthCheckResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Timestamp: time.Now().Unix(),
		Details:   make(map[string]string),
	}

	var unhealthy bool

	// pakai errgroup buat paralel check
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.postgres.GetHealth().CheckDB(ctx); err != nil {
			unhealthy = true
			response.Database = err.Error()
		} else {
			response.Database = "connected"
		}
		return nil
	})

	g.Go(func() error {
		if err := s.redisClient.GetHealth().CheckRedis(ctx); err != nil {
			unhealthy = true
			response.Redis = err.Error()
		} else {
			response.Redis = "connected"
		}
		return nil
	})

	g.Go(func() error {
		if err := s.rmq.GetHealth().CheckRabbitMQ(ctx); err != nil {
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
