package model

import "time"


type Cart struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      // Assuming users exist
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
    ID       uint `gorm:"primaryKey"`
    CartID   uint
    Cart     Cart `gorm:"foreignKey:CartID"`
    ItemID   uint
    Quantity int
    Item     Item
}
