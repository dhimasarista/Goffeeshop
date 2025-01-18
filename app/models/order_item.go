package models

import (
	"database/sql"
	"time"
)

type OrderItem struct {
	ID       sql.NullString `json:"id" gorm:"primaryKey;column:id"`
	Quantity sql.NullInt64  `gorm:"column:quantity" json:"quantity"`

	// foreign key
	OrderID   sql.NullString `json:"order_id" gorm:"column:order_id"`
	Order     Order          `gorm:"foreignKey:OrderID;references:ID"`
	ProductID sql.NullString `json:"product_id" column:"product_id"`
	Product   Product        `gorm:"foreignKey:ProductID;references:ID"`

	// Timestamp
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// Menetapkan nama tabel menjadi "order_items"
func (OrderItem) TableName() string {
	return "order_items"
}
