//go:build wireinject
// +build wireinject

package main

import (
	"Goffeeshop/app/config"
	"Goffeeshop/app/controllers"
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"

	"github.com/gofiber/fiber"
	"github.com/google/wire"
)

// InitializeApp mengonfigurasi semua dependensi aplikasi
func InitializeApp() (*fiber.App, error) {
	wire.Build(
		// Ini adalah semua dependensi yang dibutuhkan oleh main.go
		config.GormDB,
		config.NewMidtransConfig,
		repositories.NewProductRepository,
		repositories.NewOrderRepository,
		services.NewIndexService,
		services.NewOrderService,
		controllers.NewIndexController,
		controllers.NewOrderController,
		fiber.New,
	)

	return nil, nil
}
