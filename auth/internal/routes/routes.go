package routes

import (
	"auth/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.RegisterUser)
		auth.POST("/login", h.LoginUser)
		auth.POST("/refresh", h.RefreshToken)

		protected := auth.Group("/")
		protected.Use(h.RequireAuth())
		{
			protected.GET("/user", h.GetUserInfo)
			protected.PUT("/user", h.UpdateUser)
			protected.DELETE("/user", h.DeleteUser)
		}
	}

	return router
}
