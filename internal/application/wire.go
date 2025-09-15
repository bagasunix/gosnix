// internal/configs/infrastructure_wire.go
//go:build wireinject
// +build wireinject

package application

import (
	"github.com/google/wire"
	"github.com/phuslu/log"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
	"github.com/bagasunix/gosnix/internal/infrastructure/repositories"
)

func InitializeHealthHandler(
	db *gorm.DB,
	redisClient *redis.Client,
	rabbitConn *amqp091.Connection,
	logger *log.Logger,
) *handlers.HealthHandler {
	panic(wire.Build(
		// Repositories
		repositories.NewHealthRepo,
		// Services
		NewHealthService,
		// Handlers
		handlers.NewHealthHandler,
	))
}

var HealthSets = wire.NewSet(
	repositories.NewHealthRepo,
	NewHealthService,
	handlers.NewHealthHandler,
)
