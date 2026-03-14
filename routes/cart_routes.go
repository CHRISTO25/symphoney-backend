package routes

import (
	"symphoney/controllers"
	"symphoney/middleware"

	"github.com/gin-gonic/gin"
)

func CartRoutes(r *gin.Engine) {

	protected := r.Group("/cart")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/", controllers.AddToCart)

	protected.GET("/", controllers.GetCart)

	protected.PUT("/:id", controllers.UpdateCart)

	protected.DELETE("/:id", controllers.DeleteCartItem)

}
