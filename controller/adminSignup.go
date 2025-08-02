package controller

import (
	"appGO/config"
	"appGO/model"
	"appGO/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminSignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminSignupController(c *gin.Context) {
	var req AdminSignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}



	// Check for missing fields
	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Email, and Password are required"})
		return
	}

	// Check if email already exists
	var count int64
	if err := config.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password encryption failed"})
		return
	}

	// Generate OTP
	otp := utils.GenerateOTP(6)

	admin := model.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Role:       "admin",
		IsVerified: false,
		OTP:        otp,
	}

	// Save to database
	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin account"})
		return
	}

	// Send OTP via email
	if err := utils.SendEmail(req.Email, "OTP Verification for Admin", "Your admin OTP is: "+otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin account created. Please verify OTP sent to your email.",
		"admin":   admin.Email,
	})
}





