package model

import "time"

type Wishlist struct {
	ID        uint        `gorm:"primaryKey"`
	UserID    uint        `gorm:"not null"`
	Items     []WishlistItem `gorm:"foreignKey:WishlistID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WishlistItem struct {
	ID         uint      `gorm:"primaryKey"`
	WishlistID uint      `gorm:"not null"`
	ItemID     uint      `gorm:"not null"`
	Item       Item      `gorm:"foreignKey:ItemID"`
	CreatedAt  time.Time
}
