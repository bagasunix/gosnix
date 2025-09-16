package configs

import (
	"os"
	"time"

	"github.com/phuslu/log"

	"github.com/bagasunix/gosnix/pkg/configs"
)

func InitLogger(cfg *configs.Cfg) *log.Logger {
	logger := &log.Logger{
		// Level:      log.DebugLevel, // Set log level sesuai kebutuhan
		Caller:     1,            // Untuk menampilkan informasi caller
		TimeFormat: time.RFC3339, // Format waktu yang lebih mudah dibaca
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,  // Output dengan warna untuk debugging
			EndWithMessage: false, // Memastikan hanya pesan log yang ditampilkan
			QuoteString:    true,  // Quote string agar lebih mudah dibaca
		},
	}

	// Jika environment LOG_FORMAT=json, gunakan JSON writer
	if cfg.App.Environment == "json" {
		logger.Writer = &log.IOWriter{Writer: os.Stdout}
	}

	return logger
}
