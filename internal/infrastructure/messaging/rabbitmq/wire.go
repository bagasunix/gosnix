package rabbitmq

import (
	"github.com/phuslu/log"
	"github.com/rabbitmq/amqp091-go"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type RmqClient interface {
	GetHealth() repository.RabbitMQRepository
}

type repo struct {
	health repository.RabbitMQRepository
}

// GetHealth implements RmqClient.
func (r *repo) GetHealth() repository.RabbitMQRepository {
	return r.health
}

func New(logger *log.Logger, rmqConn *amqp091.Connection) RmqClient {
	return &repo{health: NewHealthRepo(logger, rmqConn)}
}
