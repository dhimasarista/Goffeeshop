package controllers

import (
	"Goffeeshop/app/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
)

type OrderController struct {
	OrderService *services.OrderService
	SocketServer *socketio.Server
}

func NewOrderController(orderService *services.OrderService, socket *socketio.Server) *OrderController {
	return &OrderController{
		OrderService: orderService,
		SocketServer: socket,
	}
}

func (controller *OrderController) GetAllOrder(ctx *fiber.Ctx) error {
	data, err := controller.OrderService.GetAllOrder()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

func (controller *OrderController) PostOrder(ctx *fiber.Ctx) error {
	postOrder, err := controller.OrderService.PostOrder(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	controller.SocketServer.BroadcastToRoom("/", "room1", "newOrder", fmt.Sprintf("New Order: %v", postOrder["order_id"]))
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"order_id":          postOrder["order_id"],
		"transaction_token": postOrder["transaction_token"],
	})
}

func (controller *OrderController) CheckPaymentStatus(ctx *fiber.Ctx) error {
	orderId := ctx.Query("id")
	response, err := controller.OrderService.CheckPaymentStatus(orderId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	controller.SocketServer.BroadcastToNamespace("/", "paymentStatus", response)
	return ctx.Status(fiber.StatusOK).JSON(response)
}
