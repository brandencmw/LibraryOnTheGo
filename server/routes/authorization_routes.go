package routes

import (
	"libraryonthego/server/authentication"
	"libraryonthego/server/middleware"

	"github.com/gin-gonic/gin"
)

func AttachAuthorizationRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authentication.LoginUser)
		authGroup.POST("/", middleware.AuthMiddleware, authentication.ValidateUser)
	}
}
