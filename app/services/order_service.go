package services

import (
	"Goffeeshop/app/models"
	"Goffeeshop/app/repositories"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type OrderService struct {
	OrderRepo   *repositories.OrderRepository
	ProductRepo *repositories.ProductRepository
}

type OrderRequest struct {
	ProductId string `json:"productId" xml:"productId" form:"productId"`
	Quantity  int    `json:"quantity" xml:"quantity" form:"quantity"`
}

func NewOrderService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository) *OrderService {
	return &OrderService{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
	} // Menggunakan OrderRepo yang bertipe *repositories.OrderRepository
}

func (os *OrderService) GetAllOrder() ([]map[string]interface{}, error) {
	orders, err := os.OrderRepo.WithOrderItem()
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *OrderService) PostOrder(ctx *fiber.Ctx) (map[string]any, error) {
	var rawBody map[string]json.RawMessage
	var orderItems []OrderRequest

	err := json.Unmarshal(ctx.Body(), &rawBody)
	if err != nil {
		return nil, err
	}

	if productRaw, exists := rawBody["products"]; exists {
		err = json.Unmarshal(productRaw, &orderItems)
		if err != nil {
			log.Println("Error OrderService.PostOrder:", err)
			return nil, err
		}
	} else {
		return nil, &fiber.Error{
			Code:    404,
			Message: "Not Found",
		}
	}

	var transactionToken string
	orderId := uuid.New().String()
	err = os.ProductRepo.DB.Transaction(func(tx *gorm.DB) error {
		var totalAmount int64
		if err = tx.Create(&models.Order{
			ID:               sql.NullString{String: orderId, Valid: true},
			Status:           sql.NullString{String: "pending", Valid: true},
			TotalAmount:      sql.NullInt64{Int64: totalAmount, Valid: true},
			TransactionToken: sql.NullString{String: "", Valid: true},
			UpdatedAt:        time.Time{},
		}).Error; err != nil {
			log.Println("Error OrderService.PostOrder:", err)
			return err
		}
		for _, item := range orderItems {
			product, err := os.ProductRepo.First(item.ProductId)
			if err != nil {
				log.Println("Error OrderService.PostOrder:", err)
				return err
			}
			amount := item.Quantity * int(product.Price.Int64)
			totalAmount += int64(amount)
			if err := tx.Create(&models.OrderItem{
				ID:        sql.NullString{String: uuid.New().String(), Valid: true},
				Quantity:  sql.NullInt64{Int64: int64(item.Quantity), Valid: true},
				Amount:    sql.NullInt64{Int64: int64(amount), Valid: true},
				OrderID:   sql.NullString{String: orderId, Valid: true},
				ProductID: sql.NullString{String: product.ID.String, Valid: true},
			}).Error; err != nil {
				log.Println("Error OrderService.PostOrder:", err)
				return err
			}
		}
		// midtrans transaksi
		req := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderId,
				GrossAmt: totalAmount,
			},
			Items: &[]midtrans.ItemDetails{},
		}
		res, _ := snap.CreateTransaction(req)
		transactionToken = res.Token
		if err := tx.Model(&models.Order{}).Where("id = ?", orderId).Update("total_amount", totalAmount).Update("transaction_token", res.Token).Error; err != nil {
			log.Println("Error OrderService.PostOrder:", err)
			return err
		}

		return nil
	})

	data := map[string]any{
		"orderId":           orderId,
		"transaction_token": transactionToken,
	}
	return data, nil
}

func (os *OrderService) CheckPaymentStatus(orderID string) (map[string]any, error) {
	url := fmt.Sprintf("https://api.sandbox.midtrans.com/v2/%s/status", orderID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Note: ambil ke config ga kebaca
	serverKey := "SB-Mid-server-zTZ2r8AhWDPPeBo7H8bWtssm"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(serverKey))
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var midtransResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&midtransResp); err != nil {
		return nil, err
	}
	err = os.OrderRepo.DB.Transaction(func(tx *gorm.DB) error {
		var status string = "pending"
		if midtransResp["transaction_status"] != nil {
			// set transaction status based on response from check transaction status
			if midtransResp["transaction_status"] == "capture" {
				if midtransResp["fraud_status"] == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
					status = "challenged"
				} else if midtransResp["fraud_status"] == "accept" {
					// TODO set transaction status on your database to 'success'
					status = "success"
				}
			} else if midtransResp["transaction_status"] == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				status = "success"
			} else if midtransResp["transaction_status"] == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
				status = "deny"
			} else if midtransResp["transaction_status"] == "cancel" || midtransResp["transaction_status"] == "expire" {
				// TODO set transaction status on your databaase to 'failure'
				status = "cancel"
			} else if midtransResp["transaction_status"] == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
				status = "waiting"
			}
		} else {
			status = "cancel"
		}

		if err := tx.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error; err != nil {
			log.Println("Error OrderService.PostOrder:", err)
			return err
		}

		return nil
	})

	return map[string]any{
		"status_code":    midtransResp["status_code"],
		"status_message": midtransResp["status_message"],
	}, nil
}
