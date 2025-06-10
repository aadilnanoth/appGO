package controller

import (
	"appGO/config"
	"appGO/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	email := c.MustGet("userEmail").(string)

	var user model.User
err := config.DB.QueryRow("SELECT id, name, email FROM users WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		fmt.Println("Query Error:", err) 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
	   "name":user.Name,
		"email": user.Email,
	})
}
