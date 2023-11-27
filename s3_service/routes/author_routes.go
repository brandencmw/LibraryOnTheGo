package routes

import (
	"s3/controllers"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine, controller *controllers.AuthorsController) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.GET("/key", controller.RetrieveObjectKey)
		authorGroup.POST("/add", controller.UploadAuthorImage)
		authorGroup.PUT("/update", controller.ReplaceAuthorImage)
		authorGroup.DELETE("/delete", controller.DeleteAuthorImage)
	}

}
