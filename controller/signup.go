package controller

import (
	"appGO/config"
	"appGO/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Name  string`json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
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

	// Check if user already exists
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists)
	if err != nil {
		log.Println("❌ QueryRow error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("❌ Password hashing error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Insert user
_, err = config.DB.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
	req.Name, req.Email, string(hashedPassword))
	if err != nil {
		log.Println("❌ Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

token, err :=utils.GenerateJWT(req.Email)
	if err != nil {
		log.Println("❌ Token generation error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}
