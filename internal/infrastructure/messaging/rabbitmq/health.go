package rabbitmq

import (
	"context"

	"github.com/phuslu/log"
	"github.com/rabbitmq/amqp091-go"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type rmqProviderHealth struct {
	logger  *log.Logger
	rmqconn *amqp091.Connection
}

func NewHealthRepo(logger *log.Logger, rmqconn *amqp091.Connection) repository.RabbitMQRepository {
	g := new(rmqProviderHealth)
	g.logger = logger
	g.rmqconn = rmqconn
	return g
}

func (g *rmqProviderHealth) CheckRabbitMQ(ctx context.Context) error {
	ch, err := g.rmqconn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return nil
}
