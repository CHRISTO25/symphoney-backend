package routes

import (
	"fmt"
)

func main() {
	fmt.Println("Hello")
}

// import (
// 	"symphoney/controllers"
// 	"symphoney/middleware"

// 	"github.com/gin-gonic/gin"
// )

// func CategoryRoutes(r *gin.Engine) {

// 	protected := r.Group("/categories")
// 	protected.Use(middleware.AuthMiddleware())

// 	protected.POST("/", controllers.CreateCategory)
// 	protected.GET("/", controllers.GetCategories)

// }
