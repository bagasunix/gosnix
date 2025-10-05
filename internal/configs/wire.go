//go:build wireinject
// +build wireinject

// internal/configs/config_wire.go
// go:build wireinject
package configs

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/phuslu/log"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/pkg/configs"
)

// Container untuk semua configs
type Configs struct {
	Cfg         *configs.Cfg
	DBPostgres  *gorm.DB
	DBMongo     *mongo.Database
	RedisClient *redis.Client
	RabbitConn  *amqp091.Connection
	Logger      *log.Logger
	FiberApp    *fiber.App
}

var ConfigSet = wire.NewSet(
	configs.InitConfig,
	InitLogger,
	InitDBPostgres,
	InitDBMongo,
	InitRedis,
	InitRabbitMQ,
	InitFiber,
	wire.Struct(new(Configs), "*"), // inject semua dependency ke struct Configs
)

func InitializeConfigs(ctx context.Context) *Configs {
	panic(wire.Build(ConfigSet))
}
