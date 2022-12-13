package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Name      string    `json:"name" xml:"name" form:"name" query:"name"`
	Password  string    `json:"password" xml:"password" form:"password" query:"password"`
	Email     string    `json:"email" xml:"email" form:"email" query:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
