package repositories

import (
	"Goffeeshop/app/models"
	"log"
	"sync"

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

/*
note:
- bukan kode optimized, gunakan query untuk performa lebih bagus
- kalo dataset nya kecil, channel & goroutine bikin overhead
*/
func (repo *OrderRepository) WithOrderItem() ([]map[string]any, error) {
	var orders []models.Order
	var orderDataFormatted []map[string]any

	err := repo.DB.Preload("OrderItems.Product").Find(&orders).Error
	if err != nil {
		log.Println("Error get all orders: ", err)
		return nil, err
	}

	// Channel untuk hasil
	ch := make(chan map[string]any, len(orders))
	wg := sync.WaitGroup{}

	// Goroutine untuk setiap order
	for _, order := range orders {
		wg.Add(1)
		go func(order models.Order) {
			defer wg.Done()
			var products []map[string]any
			for _, orderItem := range order.OrderItems {
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
				"order_id":          order.ID.String,
				"products":          products,
				"created_at":        order.CreatedAt,
				"updated_at":        order.UpdatedAt,
			}
			ch <- data
		}(order)
	}

	// Menutup channel setelah semua goroutine selesai
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Mengumpulkan hasil
	for data := range ch {
		orderDataFormatted = append(orderDataFormatted, data)
	}

	return orderDataFormatted, nil
}
