package repositories

import (
	"Goffeeshop/app/models"
	"log"

	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (repo *OrderRepository) All() ([]models.Order, error) {
	var orders []models.Order

	err := repo.DB.Find(&orders).Error
	if err != nil {
		log.Println("Error get all products: ", err)
		return nil, err
	}

	return orders, nil
}

func (repo *OrderRepository) WithOrderItem() ([]map[string]any, error) {
	var orders []models.Order
	var formattedOrders []map[string]any

	// Load orders beserta relasi OrderItems dan Product
	err := repo.DB.Debug().Preload("OrderItems.Product").Order("orders.created_at DESC").Find(&orders).Error
	if err != nil {
		log.Println("Error while fetching orders:", err)
		return nil, err
	}

	// Format data orders
	for _, order := range orders {
		// Format data products untuk setiap order
		var products []map[string]any
		for _, orderItem := range order.OrderItems {
			products = append(products, map[string]any{
				"name":     orderItem.Product.Name.String,
				"quantity": orderItem.Quantity.Int64,
				"amount":   orderItem.Amount.Int64,
			})
		}

		// Format data order
		formattedOrders = append(formattedOrders, map[string]any{
			"id":                order.ID.String,
			"status":            order.Status.String,
			"total_amount":      order.TotalAmount.Int64,
			"transaction_token": order.TransactionToken.String,
			"order_id":          order.ID.String,
			"products":          products,
			"created_at":        order.CreatedAt,
			"updated_at":        order.UpdatedAt,
		})
	}

	return formattedOrders, nil
}
