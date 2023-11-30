package routes

import (
	"s3/controllers"

	"github.com/gin-gonic/gin"
)

func AttachBookRoutes(router *gin.Engine, controller *controllers.BucketController) {
	bookGroup := router.Group("/books")
	{
		bookGroup.GET("/key", controller.RetrieveObjectKey)
		bookGroup.POST("/add", controller.UploadImage)
		bookGroup.PUT("/update", controller.ReplaceImage)
		bookGroup.DELETE("/delete", controller.DeleteImage)
	}
}
