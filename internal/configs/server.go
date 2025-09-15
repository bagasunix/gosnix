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
	cfg := InitConfig(ctx)
	container := InitializeConfigs(ctx, cfg)

	// Pastikan resource ditutup saat shutdown
	defer closeResources(container)

	// Initialize handler & setup routes
	healthHandler := application.InitializeHealthHandler(container.DB, container.RedisClient, container.RabbitConn, container.Logger)
	httpRouter.SetupRoutes(container.FiberApp, healthHandler)

	// Channel untuk menangani error atau signal
	errs := make(chan error, 1)

	// Start server & signal listener
	go runHTTP(container, cfg.Server.Port, errs)
	go listenSignal(errs)

	// Tunggu error / signal
	err := <-errs
	container.Logger.Warn().Msgf("Shutting down: %v", err)

	// Graceful shutdown
	gracefulShutdown(ctx, container)
}

// -------------------- Helper Functions --------------------

func closeResources(c *Configs) {
	sqlDB, _ := c.DB.DB()
	_ = sqlDB.Close()
	_ = c.RedisClient.Close()
	_ = c.RabbitConn.Close()
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
