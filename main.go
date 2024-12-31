package main

import (
	"Goffeeshop/app/config"
	"Goffeeshop/app/controllers"
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"
	"Goffeeshop/app/utilities"
	"fmt"
	"log"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
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

	// middlewares
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Static("/", "./public")

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

	// SocketIO
	socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) {
		fmt.Printf("Connection event 1 - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
		// Remove the user from the local clients
		delete(config.SocketIOClient, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		// Remove the user from the local clients
		delete(config.SocketIOClient, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Close event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
		fmt.Printf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	app.Get("/ws", socketio.New(func(kws *socketio.Websocket) {
		kws.Broadcast([]byte("Another User Connected"), true, socketio.TextMessage)
	}))
	// Websocket Routes using SocketIO
	app.Get("/ws/order/new", socketio.New(func(kws *socketio.Websocket) {
		userId := kws.UUID
		config.AddSocketIOClient("/ws/order/new", userId, kws)
		kws.SetAttribute("user_id", userId)

		kws.Broadcast([]byte("New Order"), true, socketio.TextMessage)

		// kws.Emit([]byte(fmt.Sprintf("Hello user: %s", userId)), socketio.TextMessage)
	}))
	app.Get("/ws/order/status", socketio.New(func(kws *socketio.Websocket) {
		userId := kws.UUID
		config.AddSocketIOClient("/ws/order/status", userId, kws)
		kws.SetAttribute("user_id", userId)

		kws.Broadcast([]byte("Order Status"), true, socketio.TextMessage)

		// kws.Emit([]byte(fmt.Sprintf("Hello user: %s", userId)), socketio.TextMessage)
	}))

	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":3000"))
}

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}
