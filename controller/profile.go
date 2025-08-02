package controller

import (
	"appGO/config"
	"appGO/model"

	"net/http"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	email := c.MustGet("userEmail").(string)

	var user model.User
 if err := config.DB.Select("name","email").Where("email = ?",email).First(&user).Error;err!=nil{
	c.JSON(http.StatusInternalServerError,gin.H{"error":"User not found"})
	return
}
	

	c.JSON(http.StatusOK, gin.H{
	   "name":user.Name,
		"email": user.Email,
	})
}
