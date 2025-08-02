package controller

import (
	"appGO/config"
	"appGO/model"
	"appGO/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)




func VerifyOTP (c *gin.Context){
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user model.User
	if err:=config.DB.Where("email=?",req.Email).First(&user).Error;
	 err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}
	user.IsVerified=true
	if err:=config.DB.Save(&user).Error;err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to Verify user"})
		return
	}
	token ,err:=utils.GenerateJWT(user.Email,user.Role)
	if err !=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":"Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"message":"OTP Verified sucessfully",
		"token":token,
	})
}


