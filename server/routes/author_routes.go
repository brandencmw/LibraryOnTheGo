package routes

import (
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine, controller *controllers.AuthorsController) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.GET("/", controller.GetAllAuthors)
		authorGroup.GET("/:authorID", controller.GetAuthor)
	}

	protectedAuthorGroup := authorGroup.Group("/auth")
	{
		protectedAuthorGroup.POST("/create", middleware.AuthMiddleware, controller.AddAuthor)
		protectedAuthorGroup.DELETE("/delete/:authorID", middleware.AuthMiddleware, controller.DeleteAuthor)
		protectedAuthorGroup.PUT("/update", middleware.AuthMiddleware, controller.UpdateAuthor)
	}
}
