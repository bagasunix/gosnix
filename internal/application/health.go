package application

import (
	"context"
	"fmt"
	"time"

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

	type checkResult struct {
		name  string
		err   error
		value string
	}

	ch := make(chan checkResult, 3) // Buffer channel untuk menghindari goroutine leak

	// Check Database
	go func() {
		err := s.postgres.GetHealth().CheckDB(ctx)
		val := "connected"
		if err != nil {
			val = err.Error()
		}
		ch <- checkResult{name: "Database", err: err, value: val}
	}()

	// Check Redis dengan pengecekan nil
	go func() {
		var err error
		var val string = "connected"

		// Pastikan redisClient tidak nil sebelum memanggil GetHealth()
		if s.redisClient == nil {
			err = fmt.Errorf("redis client is not initialized")
			val = err.Error()
		} else if s.redisClient.GetHealth() == nil {
			err = fmt.Errorf("redis health checker is not initialized")
			val = err.Error()
		} else {
			err = s.redisClient.GetHealth().CheckRedis(ctx)
			if err != nil {
				val = err.Error()
			}
		}

		ch <- checkResult{name: "Redis", err: err, value: val}
	}()

	// Check RabbitMQ
	go func() {
		err := s.rmq.GetHealth().CheckRabbitMQ(ctx)
		val := "connected"
		if err != nil {
			val = err.Error()
		}
		ch <- checkResult{name: "RabbitMQ", err: err, value: val}
	}()

	// Tunggu semua hasil
	unhealthy := false
	for i := 0; i < 3; i++ {
		res := <-ch
		switch res.name {
		case "Database":
			response.Database = res.value
		case "Redis":
			response.Redis = res.value
		case "RabbitMQ":
			response.RabbitMQ = res.value
		}
		if res.err != nil {
			unhealthy = true
		}
	}

	if unhealthy {
		response.Status = "unhealthy"
	}

	return response, nil
}
