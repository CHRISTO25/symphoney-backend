package routes

import (
	"symphoney/controllers"
	"symphoney/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {

	protected := r.Group("/products")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/", controllers.CreateProduct)

	protected.GET("/", controllers.GetProducts)

	protected.GET("/:id", controllers.GetProductByID)

	protected.PUT("/:id", controllers.UpdateProduct)

	protected.DELETE("/:id", controllers.DeleteProduct)

	protected.POST("/:id/images", controllers.AddProductImage)

}
