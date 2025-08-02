package main

import (
	"appGO/config"
	"appGO/route"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	
	r := gin.Default()
	route.SetupAdminRoutes(r)
	route.SetupRoutes(r)
    r.Run("0.0.0.0:8080")
}
