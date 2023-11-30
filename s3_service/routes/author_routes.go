package routes

import (
	"s3/controllers"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine, controller *controllers.BucketController) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.GET("/key", controller.RetrieveObjectKey)
		authorGroup.POST("/add", controller.UploadImage)
		authorGroup.PUT("/update", controller.ReplaceImage)
		authorGroup.DELETE("/delete", controller.DeleteImage)
	}

}
