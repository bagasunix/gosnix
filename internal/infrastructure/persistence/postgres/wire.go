package postgres

import (
	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type Repositories interface {
	GetCustomer() repository.CustomerRepository
	GetHealth() repository.PostgresRepository
}

type repo struct {
	customer repository.CustomerRepository
	health   repository.PostgresRepository
}

// GetCustomer implements Repositories.
func (r *repo) GetCustomer() repository.CustomerRepository {
	return r.customer
}

// GetHealth implements Repositories.
func (r *repo) GetHealth() repository.PostgresRepository {
	return r.health
}

func New(logger *log.Logger, db *gorm.DB) Repositories {
	return &repo{
		customer: NewGormCustomer(logger, db),
		health:   NewHealthRepo(logger, db),
	}
}
