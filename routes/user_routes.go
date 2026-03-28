package routes

import (
	"github.com/gin-gonic/gin"
	"symphoney/controllers"
	"symphoney/middleware"
)

func UserRoutes(r *gin.Engine) {
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	protected := r.Group("/profile")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/", controllers.Profile)

	protected.GET("/users", controllers.GetAllUsers)
	protected.PUT("/users/block/:id", controllers.BlockUser)
	protected.PUT("/users/unblock/:id", controllers.UnblockUser)
}
