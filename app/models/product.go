package models

import (
	"database/sql"
)

type Product struct {
	ID    sql.NullString `json:"id" gorm:"primaryKey;column:id"`
	Name  sql.NullString `json:"name" gorm:"column:name"`
	Price sql.NullInt64  `gorm:"column:price" json:"price"`

	// Timestamp
	// CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	// UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	// DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
