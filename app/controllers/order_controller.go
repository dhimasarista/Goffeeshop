package controllers

import (
	"Goffeeshop/app/services"

	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	OrderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{OrderService: orderService}
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
	data, err := controller.OrderService.PostOrder(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"order_id":          data["order_id"],
		"transaction_token": data["transaction_token"],
	})
}

func (controller *OrderController) CheckPaymentStatus(ctx *fiber.Ctx) error {
	orderId := ctx.Query("id")
	response, err := controller.OrderService.CheckPaymentStatus(orderId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
