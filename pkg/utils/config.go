package utils

import (
	"fmt"
	"time"
)

type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	User            string
	Password        string
	DatabaseName    string
	SSLMode         string
	Timezone        string
	MaxOpenConns    int           // tambahan: max open connections
	MaxIdleConns    int           // tambahan: idle timeout
	ConnMaxLifetime time.Duration // tambahan: lifetime timeout
	ConnMaxIdleTime time.Duration // tambahan: idle timeout
}

func (d *DBConfig) GetDSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/", d.Driver, d.User, d.Password, d.Host, d.Port)
}
