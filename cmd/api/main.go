package main

import (
	"github.com/bagasunix/gosnix/internal/configs"
)

// @title Gosnix API
// @version 1.0
// @description API untuk sistem Gosnix
// @termsOfService http://swagger.io/terms/

// @contact.name Developer Support
// @contact.email support@gosnix.local

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	configs.Run()
}
