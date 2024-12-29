package controllers

import (
	"Goffeeshop/app/repositories"
	"Goffeeshop/app/services"

	"github.com/gofiber/fiber/v2"
)

type IndexController struct {
	Indexservice *services.IndexService
	ProductRepo  *repositories.ProductRepository
}

func NewIndexController(indexService *services.IndexService, productRepo *repositories.ProductRepository) *IndexController {
	return &IndexController{
		Indexservice: indexService,
		ProductRepo:  productRepo,
	}
}

func (controller *IndexController) Index(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"Title": "Goffeeshop",
	}, "layouts/main")
}

func (controller *IndexController) ListOrder(ctx *fiber.Ctx) error {
	return ctx.Render("list_order", fiber.Map{})
}
func (controller *IndexController) NewOrder(ctx *fiber.Ctx) error {
	products, err := controller.ProductRepo.All()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	var simplifiedData []map[string]any
	for _, product := range products {
		data := map[string]any{
			"id":    product.ID.String,
			"name":  product.Name.String,
			"price": product.Price.Int64,
		}

		simplifiedData = append(simplifiedData, data)
	}

	return ctx.Render("new_order", fiber.Map{
		"products": simplifiedData,
	})
}