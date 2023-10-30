package routes

import (
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.GET("/", controllers.GetAllAuthors)
		authorGroup.GET("/:authorID", controllers.GetAuthor)
	}

	protectedAuthorGroup := authorGroup.Group("/auth")
	{
		protectedAuthorGroup.POST("/create", middleware.AuthMiddleware, controllers.AddAuthor)
		protectedAuthorGroup.DELETE("/delete/:authorID", middleware.AuthMiddleware, controllers.DeleteAuthor)
		protectedAuthorGroup.PUT("/update", middleware.AuthMiddleware, controllers.UpdateAuthor)
	}
}
