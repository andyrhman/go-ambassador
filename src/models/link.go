package models

import "github.com/google/uuid"

type Link struct {
	Model
	Code     string    `json:"code"`
	UserId   uuid.UUID `json:"user_id"`
	User     User      `json:"user" gorm:"foreignKey:UserId"`
	Products []Product `json:"products" gorm:"many2many:link_products"`
	Orders   []Order   `json:"orders,omitempty" gorm:"-"` // Remove the null data using omitempty
}
