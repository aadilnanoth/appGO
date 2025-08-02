package controller

import (
	"appGO/config"
	"appGO/model"
	"appGO/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupController(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	var count int64
	if err := config.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		log.Println("❌ DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Password hashing error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	otp := utils.GenerateOTP(6)

	user := model.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   string(hashedPassword),
		OTP:        otp,
		IsVerified: false,
		Role:       "user", // 👈 Default role is set here
	}

	if err := config.DB.Create(&user).Error; err != nil {
		log.Println("❌ Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	if err := utils.SendEmail(req.Email, "OTP Verification", "Your OTP is: "+otp); err != nil {
		log.Println("❌ Email sending error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP email"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully. Please verify OTP sent to your email."})
}
