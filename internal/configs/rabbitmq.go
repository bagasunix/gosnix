package configs

import (
	"context"
	"strconv"
	"time"

	"github.com/phuslu/log"
	"github.com/rabbitmq/amqp091-go"
	"github.com/streadway/amqp"

	"github.com/bagasunix/gosnix/pkg/errors"
	"github.com/bagasunix/gosnix/pkg/utils"
)

func InitRabbitMQ(ctx context.Context, cfg *Cfg, logger *log.Logger) *amqp091.Connection {
	CfgBuild := &utils.DBConfig{
		Driver:   cfg.RabbitMQ.Driver,
		Host:     cfg.RabbitMQ.Host,
		Port:     strconv.Itoa(cfg.RabbitMQ.Port),
		User:     cfg.RabbitMQ.User,
		Password: cfg.RabbitMQ.Password,
	}
	return NewPRabbitMQDB(ctx, CfgBuild, logger)
}

// InitRabbitMQ initializes and returns a RabbitMQ connection
func NewPRabbitMQDB(ctx context.Context, cfg *utils.DBConfig, logger *log.Logger) *amqp091.Connection {
	// Membuat DSN RabbitMQ
	conn, err := amqp091.DialConfig(cfg.GetDSN(), amqp091.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Dial:      amqp.DefaultDial(30 * time.Second),
	})
	errors.HandlerWithOSExit(logger, err, "init", "RabbitMQ", "config", cfg.GetDSN())

	logger.Info().Msg("Connected to RabbitMQ")
	return conn
}
