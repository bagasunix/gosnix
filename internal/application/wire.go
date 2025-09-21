//go:build wireinject
// +build wireinject

// internal/appliaction/appliaction_wire.go
// go:build wireinject
package application

import (
	"github.com/google/wire"
	"github.com/phuslu/log"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/infrastructure/http/handlers"
	"github.com/bagasunix/gosnix/internal/infrastructure/messaging/rabbitmq"
	"github.com/bagasunix/gosnix/internal/infrastructure/persistence/postgres"
	redisC "github.com/bagasunix/gosnix/internal/infrastructure/persistence/redis_client"
	"github.com/bagasunix/gosnix/pkg/configs"
)

type HandlerContainer struct {
	Health   *handlers.HealthHandler
	Customer *handlers.CustomerHandler
	Repo     postgres.Repositories
}

// Entry point buat bikin semua handler
func InitializeServiceHandler(
	db *gorm.DB,
	redis *redis.Client,
	rabbitConn *amqp091.Connection,
	logger *log.Logger,
	cfg *configs.Cfg,
) *HandlerContainer {
	wire.Build(
		HealthSet,
		wire.Struct(new(HandlerContainer), "*"), // <- otomatis isi struct
	)
	return nil
}

var HealthSet = wire.NewSet(
	// Infrastructure
	postgres.New,
	redisC.New,
	rabbitmq.New,

	// Services
	NewHealthService,
	NewCustomerService,

	// Handlers
	handlers.NewHealthHandler,
	handlers.NewCustomerHandler,
)
