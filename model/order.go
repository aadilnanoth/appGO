package model

import "time"

// ----------------- DB Models -----------------

type Order struct {
	ID              uint        `gorm:"primaryKey"`
	UserID          uint        `gorm:"not null"`
	User            User        `gorm:"foreignKey:UserID"`
	TotalPrice      float64     `gorm:"not null"`
	PaymentMethod   string      `gorm:"size:50;default:COD"`
	Status          string      `gorm:"size:50;default:Pending"`
	ShippingAddress string      `gorm:"not null"`
	Phone           string      `gorm:"size:20"`
	CreatedAt       time.Time
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID       uint    `gorm:"primaryKey"`
	OrderID  uint    `gorm:"not null"`
	ItemID   uint    `gorm:"not null"`
	Item     Item    `gorm:"foreignKey:ItemID"`
	Quantity int     `gorm:"not null"`
	Price    float64 `gorm:"not null"`

	Order `gorm:"-"` // prevent circular JSON
}

// ----------------- Response Structs -----------------

type OrderResponse struct {
	ID              uint                `json:"id"`
	UserID          uint                `json:"user_id"`
	User            UserSummary         `json:"user"`
	TotalPrice      float64             `json:"total_price"`
	PaymentMethod   string              `json:"payment_method"`
	Status          string              `json:"status"`
	ShippingAddress string              `json:"shipping_address"`
	Phone           string              `json:"phone"`
	CreatedAt       time.Time           `json:"created_at"`
	OrderItems      []OrderItemResponse `json:"order_items"`
}

type OrderItemResponse struct {
	ID       uint       `json:"id"`
	OrderID  uint       `json:"order_id"`
	ItemID   uint       `json:"item_id"`
	Item     ItemSimple `json:"item"`
	Quantity int        `json:"quantity"`
	Price    float64    `json:"price"`
}

type ItemSimple struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	ImageURL    string         `json:"image_url"`
	Category    CategorySimple `json:"category"`
}

type CategorySimple struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
