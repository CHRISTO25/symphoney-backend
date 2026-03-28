package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"symphoney/config"
	"symphoney/routes"
)

func main() {
	config.ConnectDatabase()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all (dev only)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	routes.UserRoutes(r)
	routes.CategoryRoutes(r)
	routes.ProductRoutes(r)
	routes.CartRoutes(r)
	r.Static("/uploads", "./uploads")
	r.Run(":8081")
}
