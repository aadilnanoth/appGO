package main

import (
	"appGO/config"

	"appGO/route"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	config.InitDB()
	r := gin.Default()
	route.SetupRoutes(r)
	r.Run()
}
