// internal/configs/config_wire.go
// go:build wireinject
//go:build wireinject
// +build wireinject

package configs

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/phuslu/log"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Container untuk semua configs
type Configs struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	RabbitConn  *amqp091.Connection
	Logger      *log.Logger
	FiberApp    *fiber.App
}

var ConfigSet = wire.NewSet(
	InitLogger,
	InitDB,
	InitRedis,
	InitRabbitMQ,
	InitFiber,
	wire.Struct(new(Configs), "*"), // inject semua dependency ke struct Configs
)

func InitializeConfigs(ctx context.Context, cfg *Cfg) *Configs {
	panic(wire.Build(ConfigSet))
}
