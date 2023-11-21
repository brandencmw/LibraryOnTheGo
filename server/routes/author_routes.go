package routes

import (
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine, controller *controllers.AuthorsController) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.GET("", controller.GetAuthor)
	}

	protectedAuthorGroup := authorGroup.Group("/auth", middleware.AuthMiddleware)
	{
		protectedAuthorGroup.POST("/create", controller.AddAuthor)
		protectedAuthorGroup.DELETE("/delete", controller.DeleteAuthor)
		protectedAuthorGroup.PUT("/update", controller.UpdateAuthor)
	}
}
