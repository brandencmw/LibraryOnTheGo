package routes

import (
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"

	"github.com/gin-gonic/gin"
)

func AttachBookRoutes(router *gin.Engine, controller *controllers.BooksController) {
	authorGroup := router.Group("/books")
	{
		authorGroup.GET("", controller.GetBook)
	}

	protectedAuthorGroup := authorGroup.Group("/auth", middleware.AuthMiddleware)
	{
		protectedAuthorGroup.POST("/create", controller.AddBook)
		protectedAuthorGroup.DELETE("/delete", controller.DeleteBook)
		protectedAuthorGroup.PUT("/update", controller.UpdateBook)
	}
}
