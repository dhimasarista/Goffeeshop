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

// bukan kode optimized, best practice gunakan query langsung
func (repo *OrderRepository) WithOrderItem() ([]map[string]any, error) {
	var orders []models.Order

	// preload == eager loading
	err := repo.DB.Preload("OrderItems.Product").Model(&models.Order{}).Find(&orders).Error
	if err != nil {
		log.Println("Error get all orders: ", err)
		return nil, err
	}

	var newOrdersData []map[string]any

	for _, order := range orders {
		var orderId string
		var products []map[string]any
		for _, orderItem := range order.OrderItems {
			orderId = orderItem.OrderID.String
			product := map[string]any{
				"name":     orderItem.Product.Name.String,
				"quantity": orderItem.Quantity.Int64,
				"amount":   orderItem.Amount.Int64,
			}
			products = append(products, product)
		}

		data := map[string]any{
			"id":                order.ID.String,
			"status":            order.Status.String,
			"total_amount":      order.TotalAmount.Int64,
			"transaction_token": order.TransactionToken.String,
			"order_id":          orderId,
			"products":          products,
			"created_at":        order.CreatedAt,
			"updated_at":        order.UpdatedAt,
		}
		newOrdersData = append(newOrdersData, data)
	}

	return newOrdersData, nil
}
