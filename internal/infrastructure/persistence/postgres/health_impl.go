package postgres

import (
	"context"

	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/bagasunix/gosnix/internal/domain/repository"
)

type gormProviderHealth struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewHealthRepo(logger *log.Logger, db *gorm.DB) repository.PostgresRepository {
	g := new(gormProviderHealth)
	g.db = db
	g.logger = logger
	return g
}

func (g *gormProviderHealth) CheckDB(ctx context.Context) error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
