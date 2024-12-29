package main

import (
	"Goffeeshop/app/config"
	"Goffeeshop/app/controllers"
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"
	"Goffeeshop/app/utilities"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/mustache/v2"
)

func main() {
	utilities.ClearScreen()
	// DB Config
	db := config.GormDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}

	engine := mustache.New("./views", ".mustache")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/") // Handling for nothing routes
		},
	})
	app.Static("/", "./public") // Middleware untuk menyediakan file static

	// init midtrans
	config.NewMidtransConfig()

	// Init Repo
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Init Service
	indexService := services.NewIndexService(productRepo)
	orderService := services.NewOrderService(orderRepo, productRepo)

	// Init Controller
	indexController := controllers.NewIndexController(indexService, productRepo)
	orderController := controllers.NewOrderController(orderService)

	// Routes Web
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/order")
	})
	app.Get("/order", indexController.Index)
	app.Get("/order/list", indexController.ListOrder)
	app.Get("/order/new", indexController.NewOrder)

	// Routes Api
	app.Route("/api", func(api fiber.Router) {
		api.Post("/order", orderController.PostOrder)
		api.Get("/order/list", orderController.GetAllOrder).Name("")
		api.Get("/order/check-status", orderController.CheckPaymentStatus)
	})

	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":3000"))
}
