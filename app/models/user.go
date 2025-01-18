package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        sql.NullString `json:"id" gorm:"primaryKey;column:id"`
	Username  sql.NullString `json:"username" gorm:"column:username"`
	Password  sql.NullString `json:"password" gorm:"column:password"`
	Email     sql.NullString `json:"email" gorm:"column:email"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index"`
}
