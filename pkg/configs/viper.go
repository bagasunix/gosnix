package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/bagasunix/gosnix/pkg/errors"
)

type Cfg struct {
	App struct {
		Name        string `mapstructure:"name"`
		Version     string `mapstructure:"version"`
		Environment string `mapstructure:"environment"`
		TimeZone    string `mapstructure:"time_zone"`
	} `mapstructure:"app"`

	Server struct {
		Port        int    `mapstructure:"port"`
		Version     string `mapstructure:"version"`
		RateLimiter struct {
			Enabled  bool          `mapstructure:"enabled"`
			Limit    int           `mapstructure:"limit"`
			Duration time.Duration `mapstructure:"duration"`
		} `mapstructure:"rate_limiter"`
		MailJet struct {
			ApiKey   string `mapstructure:"api_key"`
			ScretKey string `mapstructure:"scret_key"`
			HostName string `mapstructure:"hostname"`
			Port     int    `mapstructure:"port"`
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"mailjet"`
		Token struct {
			JWTKey       string `mapstructure:"jwt_key"`
			SignatureKey string `mapstructure:"signature_key"`
		} `mapstructure:"token"`
	} `mapstructure:"server"`

	Database struct {
		Driver        string        `mapstructure:"driver"`
		Host          string        `mapstructure:"host"`
		Port          int           `mapstructure:"port"`
		User          string        `mapstructure:"user"`
		Password      string        `mapstructure:"password"`
		DBName        string        `mapstructure:"dbname"`
		SSLMode       string        `mapstructure:"sslmode"`
		MaxConnection int           `mapstructure:"max_connection"`
		MaxIdleConns  int           `mapstructure:"max_idle"`
		MaxLifeTime   time.Duration `mapstructure:"max_life_time"`
		MaxIdleTime   time.Duration `mapstructure:"max_idle_time"`
		MigrationPath string        `mapstructure:"migration_path"`
	} `mapstructure:"database"`

	RabbitMQ struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"rabbitmq"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		Type     string `mapstructure:"type"`
	} `mapstructure:"redis"`

	Logging struct {
		Level  int    `mapstructure:"level"`
		Format string `mapstructure:"format"`
	} `mapstructure:"logging"`
}

func LoadCfg(ctx context.Context) (*Cfg, error) {
	// Here you can implement loading configuration from file, environment variables, etc.
	// For simplicity, we'll just return nil for now.
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".") // root project
	viper.AddConfigPath("./")
	viper.AddConfigPath("../../") // jaga2 kalau ada perbedaan path

	viper.AutomaticEnv()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Cfg
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.Database.Driver == "" || config.Database.Host == "" || config.Database.Port == 0 || config.Database.User == "" || config.Database.Password == "" || config.Database.DBName == "" {
		return nil, errors.CustomError("database configuration is missing")
	}

	if config.RabbitMQ.Driver == "" || config.RabbitMQ.Host == "" || config.RabbitMQ.Port == 0 || config.RabbitMQ.User == "" || config.RabbitMQ.Password == "" {
		return nil, errors.CustomError("rabbitmq configuration is missing")
	}

	if config.Redis.Host == "" || config.Redis.Port == "" {
		return nil, errors.CustomError("redis configuration is missing")
	}

	return &config, nil
}

// Fungsi untuk menginisialisasi konfigurasi
func InitConfig(ctx context.Context) *Cfg {
	config, err := LoadCfg(ctx)
	if err != nil {
		fmt.Println("cannot load config")
		os.Exit(1)
	}

	fmt.Println("Configuration loaded successfully")
	return config
}
