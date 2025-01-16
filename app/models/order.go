package models

import (
	"database/sql"
	"time"
)

type Order struct {
	ID               sql.NullString `json:"id" gorm:"primaryKey;column:id"`
	Status           sql.NullString `json:"status" gorm:"column:status"`
	TotalAmount      sql.NullInt64  `gorm:"column:total_amount" json:"total_amount"`
	TransactionToken sql.NullString `json:"transaction_token" gorm:"column:transaction_token"`

	// One To Many
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID"`

	// Timestamp
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
