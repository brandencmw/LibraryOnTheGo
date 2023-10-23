package routes

import (
	"libraryonthego/server/controllers"
	"libraryonthego/server/middleware"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.POST("/create", middleware.AuthMiddleware, controllers.AddAuthor)
		authorGroup.DELETE("/delete/:authorID", middleware.AuthMiddleware, controllers.DeleteAuthor)
		authorGroup.PUT("/update", middleware.AuthMiddleware, controllers.UpdateAuthor)
		authorGroup.GET("/", controllers.GetAllAuthors)
		authorGroup.GET("/:authorID", controllers.GetAuthor)
	}
}
