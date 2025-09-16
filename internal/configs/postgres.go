package configs

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migPostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/bagasunix/gosnix/pkg/configs"
	"github.com/bagasunix/gosnix/pkg/errors"
	"github.com/bagasunix/gosnix/pkg/utils"
)

func InitDB(ctx context.Context, cfg *configs.Cfg, logger *log.Logger) *gorm.DB {
	CfgBuild := &utils.DBConfig{
		Driver:          cfg.Database.Driver,
		Host:            cfg.Database.Host,
		Port:            strconv.Itoa(cfg.Database.Port),
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DatabaseName:    cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxConnection,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.MaxLifeTime,
		ConnMaxIdleTime: cfg.Database.MaxIdleTime,
		Timezone:        cfg.App.TimeZone,
	}
	return NewPostgresDB(ctx, CfgBuild, cfg.Database.MigrationPath, logger)
}

func NewPostgresDB(ctx context.Context, cfg *utils.DBConfig, migrationPath string, logger *log.Logger) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()+cfg.DatabaseName+"?sslmode="+cfg.SSLMode), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	errors.HandlerWithOSExit(logger, err, "init", "database", "config", cfg.GetDSN())

	// Konfigurasi koneksi DB
	errors.HandlerWithOSExit(logger, db.WithContext(ctx).Use(dbresolver.Register(dbresolver.Config{}).
		SetMaxOpenConns(cfg.MaxOpenConns).
		SetMaxIdleConns(cfg.MaxIdleConns).
		SetConnMaxLifetime(cfg.ConnMaxLifetime).
		SetConnMaxIdleTime(cfg.ConnMaxIdleTime)),
		"db_resolver")

	sqlDB, _ := db.DB()

	// Jalankan migrasi di fungsi terpisah, path fleksibel
	// runMigrations(sqlDB, migrationPath, logger)

	// Verifikasi koneksi
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		errors.HandlerWithOSExit(logger, err, "init", "database", "ping", "")
	}

	return db
}

// runMigrations akan mengurus semua proses migrasi database
func runMigrations(sqlDB *sql.DB, migrationsPath string, logger *log.Logger) {
	driver, err := migPostgres.WithInstance(sqlDB, &migPostgres.Config{})
	errors.HandlerWithOSExit(logger, err, "failed to initialize postgres driver")

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	errors.HandlerWithOSExit(logger, err, "failed to create migration instance")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		errors.HandlerWithOSExit(logger, err, "Failed to run migrations.")
	}

	fmt.Println("Migration successful")
}
