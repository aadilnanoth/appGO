package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	email := c.MustGet("userEmail").(string) 
	c.JSON(http.StatusOK, gin.H{"message": "Welcome " + email})
}