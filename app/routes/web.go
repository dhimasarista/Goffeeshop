package routes

import (
	"Goffeeshop/app/controllers"
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"

	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"gorm.io/gorm"
)

func Web(app *fiber.App, db *gorm.DB, server *socketio.Server) {
	// Init Repo
	productRepo := repositories.NewProductRepository(db)

	// Init Service
	indexService := services.NewIndexService(productRepo)

	// Init Controller
	indexController := controllers.NewIndexController(indexService, productRepo)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/order")
	})
	app.Get("/order", indexController.Index)
	app.Get("/order/list", indexController.ListOrder)
	app.Get("/order/new", indexController.NewOrder)
}
