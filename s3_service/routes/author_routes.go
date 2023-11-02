package routes

import (
	"s3/controllers"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine, controller *controllers.AuthorsController) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.POST("/add-author-image", controller.UploadAuthorImage)
		authorGroup.GET("/:object-key", controllers.RetrieveAuthorImage)
	}

}
