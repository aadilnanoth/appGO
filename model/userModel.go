package model

import "time"

// User represents a user account in the system.
type User struct {
    ID         uint      `gorm:"primaryKey"`                              // Unique identifier for the user.
    Email      string    `gorm:"uniqueIndex;not null"`                    // Email address of the user, must be unique.
    Password   string    `gorm:"not null"`                                // Hashed password of the user.
    CreatedAt  time.Time `gorm:"autoCreateTime"`                          // Timestamp when the user was created.
    Name       string    `gorm:"type:varchar(100);not null"`              // Full name of the user.
    OTP        string    `gorm:"size:6"`                                  // One-Time Password for verification.
    IsVerified bool      `gorm:"default:false"`                           // Indicates if the user has verified their account.
    Role       string    `gorm:"type:varchar(20);default:'user'"`         // Role of the user: 'user' or 'admin'.
}
