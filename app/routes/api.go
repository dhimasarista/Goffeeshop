package routes

import (
	"Goffeeshop/app/controllers"
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"

	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"
)

func ApiRoutes(app *fiber.App, db *gorm.DB, server *socketio.Server) {
	// Init Repo
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Init Service
	indexService := services.NewIndexService(productRepo)
	orderService := services.NewOrderService(orderRepo, productRepo)

	// Init Controller
	indexController := controllers.NewIndexController(indexService, productRepo)
	orderController := controllers.NewOrderController(orderService, server)
	app.Route("/api", func(api fiber.Router) {
		api.Post("/order", orderController.PostOrder)
		api.Get("/order/list", orderController.GetAllOrder)
		api.Get("/order/check-status", orderController.CheckPaymentStatus)
		api.Get("/product/list", indexController.ProductList)
	})
}
