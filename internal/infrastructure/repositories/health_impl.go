package repositories

import (
	"context"

	"github.com/phuslu/log"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/health"
)

type gormProvider struct {
	db      *gorm.DB
	logger  *log.Logger
	client  *redis.Client
	rmqconn *amqp091.Connection
}

func NewHealthRepo(logger *log.Logger, db *gorm.DB, client *redis.Client, rmqconn *amqp091.Connection) health.Repository {
	g := new(gormProvider)
	g.db = db
	g.logger = logger
	g.client = client
	g.rmqconn = rmqconn
	return g
}

func (g *gormProvider) CheckDB(ctx context.Context) error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (g *gormProvider) CheckRedis(ctx context.Context) error {
	return g.client.Ping(ctx).Err()
}

func (g *gormProvider) CheckRabbitMQ(ctx context.Context) error {
	ch, err := g.rmqconn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return nil
}
