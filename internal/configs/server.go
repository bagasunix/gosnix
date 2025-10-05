package configs

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bagasunix/gosnix/internal/application"
	httpRouter "github.com/bagasunix/gosnix/internal/infrastructure/http/router"
)

func Run() {
	// Context global untuk semua dependency
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config & initialize dependencies
	container := InitializeConfigs(ctx)

	// Pastikan resource ditutup saat shutdown
	defer closeResources(container)

	// Initialize handler & setup routes
	serviceHandler := application.InitializeServiceHandler(container.DBPostgres, container.RedisClient, container.RabbitConn, container.Logger, container.Cfg)
	httpRouter.SetupRoutes(container.FiberApp, container.Cfg, serviceHandler.Health, serviceHandler.Customer)

	// Channel untuk menangani error atau signal
	errs := make(chan error, 1)

	// Start server & signal listener
	go runHTTP(container, container.Cfg.Server.Port, errs)
	go listenSignal(errs)

	// Tunggu error / signal
	err := <-errs
	container.Logger.Warn().Msgf("Shutting down: %v", err)

	// Graceful shutdown
	gracefulShutdown(ctx, container)
}

// -------------------- Helper Functions --------------------

func closeResources(c *Configs) {
	// Tutup DBPostgres
	if c.DBPostgres != nil {
		sqlDB, err := c.DBPostgres.DB()
		if err == nil && sqlDB != nil {
			_ = sqlDB.Close()
		}
	}

	// Tutup Redis
	if c.RedisClient != nil {
		_ = c.RedisClient.Close()
	}

	// Tutup RabbitMQ
	if c.RabbitConn != nil {
		_ = c.RabbitConn.Close()
	}
}

func runHTTP(app *Configs, port int, errs chan error) {
	if err := app.FiberApp.Listen(":" + strconv.Itoa(port)); err != nil {
		errs <- fmt.Errorf("failed to start server: %w", err)
	}
}

func listenSignal(errs chan error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	errs <- fmt.Errorf("received signal: %s", <-sig)
}

func gracefulShutdown(ctx context.Context, c *Configs) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := c.FiberApp.ShutdownWithContext(timeoutCtx); err != nil {
		c.Logger.Fatal().Err(err).Msg("Server shutdown failed")
	}
	c.Logger.Info().Msg("Server shutdown gracefully")
}
