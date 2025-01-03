package main

import (
	"Goffeeshop/app/config"
	"Goffeeshop/app/controllers"
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"
	"Goffeeshop/app/utilities"
	"log"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/mustache/v2"

	socketio "github.com/googollee/go-socket.io"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func main() {
	utilities.ClearScreen()
	// DB Config
	db := config.GormDB()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}
	// SocketIO
	server := socketio.NewServer(nil) // socketio
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("Client connected:", s.ID())
		s.Join("room1") // Contoh join ke sebuah room
		return nil
	})
	// Event handler untuk menerima pesan
	server.OnEvent("/", "message", func(s socketio.Conn, msg string) {
		log.Printf("Received message: %s from %s", msg, s.ID())
		s.Emit("reply", "Message received: "+msg)
	})
	// Event handler untuk menerima pesan
	server.OnEvent("/", "newOrder", func(s socketio.Conn, msg string) {
		s.Emit("newOrder", "New Order: "+msg)
	})

	// Event handler untuk disconnect
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Client disconnected:", s.ID(), "Reason:", reason)
	})
	// Goroutine untuk menjalankan Socket.IO server
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()
	defer server.Close()

	// Fiber App
	engine := mustache.New("./views", ".mustache")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Redirect("/") // Handling for nothing routes
		},
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// middlewares
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
	orderController := controllers.NewOrderController(orderService, server)

	// Routes Web
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/order")
	})
	app.Get("/order", indexController.Index)
	app.Get("/order/list", indexController.ListOrder)
	app.Get("/order/new", indexController.NewOrder)

	app.Get("/call", func(ctx *fiber.Ctx) error {
		server.BroadcastToRoom("/", "room1", "newOrder", "New Order")
		return ctx.JSON(fiber.Map{})
	})

	// Routes Api
	app.Route("/api", func(api fiber.Router) {
		api.Post("/order", orderController.PostOrder)
		api.Get("/order/list", orderController.GetAllOrder).Name("")
		api.Get("/order/check-status", orderController.CheckPaymentStatus)
	})

	// Integrasi Socket.IO dengan Fiber
	app.All("/socket.io/*", func(c *fiber.Ctx) error {
		// Gunakan fasthttpadaptor untuk menjembatani handler
		fasthttpadaptor.NewFastHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			server.ServeHTTP(w, r)
		}))(c.Context())
		return nil
	})

	app.Get("/metrics", monitor.New())

	log.Fatal(app.Listen(":3000"))
}
