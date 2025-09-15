package health

import "context"

type PostgresRepository interface {
	CheckDB(ctx context.Context) error
}

type RedisRepository interface {
	CheckRedis(ctx context.Context) error
}

type RabbitMQRepository interface {
	CheckRabbitMQ(ctx context.Context) error
}

// Composite interface yang menggabungkan semua repository
type Repository interface {
	PostgresRepository
	RedisRepository
	RabbitMQRepository
}
