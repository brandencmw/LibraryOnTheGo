package routes

import (
	"s3/controllers"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine, controller *controllers.AuthorsController) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.POST("/add", controller.UploadAuthorImage)
		authorGroup.GET("", controller.RetrieveAuthorImage)
	}

}
