package configs

import (
	"context"
	"strconv"
	"time"

	"github.com/phuslu/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/bagasunix/gosnix/pkg/configs"
	"github.com/bagasunix/gosnix/pkg/errors"
	"github.com/bagasunix/gosnix/pkg/utils"
)

func InitDBMongo(ctx context.Context, cfg *configs.Cfg, logger *log.Logger) *mongo.Database {
	CfgBuild := &utils.DBConfig{
		Driver:          cfg.Database.MongoDB.Driver,
		Host:            cfg.Database.MongoDB.Host,
		Port:            strconv.Itoa(cfg.Database.MongoDB.Port),
		User:            cfg.Database.MongoDB.User,
		Password:        cfg.Database.MongoDB.Password,
		DatabaseName:    cfg.Database.MongoDB.DBName,
		SSLMode:         cfg.Database.MongoDB.SSLMode,
		MaxOpenConns:    cfg.Database.MongoDB.MaxPoolSize,
		MaxIdleConns:    cfg.Database.MongoDB.MinPoolSize,
		ConnMaxIdleTime: cfg.Database.MongoDB.MaxIdleTime,
		Timezone:        cfg.App.TimeZone,
	}
	return NewMongoDB(ctx, CfgBuild, logger)
}

func NewMongoDB(ctx context.Context, cfg *utils.DBConfig, logger *log.Logger) *mongo.Database {
	// Implementasi koneksi MongoDB di sini
	// Gunakan cfg untuk mendapatkan detail koneksi
	clientOpts := options.Client().
		ApplyURI(cfg.GetDSN()).
		SetMaxPoolSize(uint64(cfg.MaxOpenConns)).
		SetMinPoolSize(uint64(cfg.MaxIdleConns)).
		SetMaxConnIdleTime(cfg.ConnMaxIdleTime)

	client, err := mongo.Connect(clientOpts)
	errors.HandlerWithOSExit(logger, err, "init", "database", "mongo", cfg.GetDSN())
	// Ping DB
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = client.Ping(ctxPing, nil)
	errors.HandlerWithOSExit(logger, err, "ping", "database", "mongo", cfg.GetDSN())

	logger.Info().Msg("Connected to MongoDB")

	return client.Database(cfg.DatabaseName)
}
