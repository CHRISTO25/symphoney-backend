package main

import (
	"github.com/gin-gonic/gin"
	"symphoney/config"
	"symphoney/routes"
)

func main() {
	config.ConnectDatabase()
	r := gin.Default()
	routes.UserRoutes(r)
	routes.CategoryRoutes(r)
	routes.ProductRoutes(r)
	routes.CartRoutes(r)
	r.Run(":8080")
}
