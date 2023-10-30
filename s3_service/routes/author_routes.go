package routes

import (
	"s3/controllers"

	"github.com/gin-gonic/gin"
)

func AttachAuthorRoutes(router *gin.Engine) {
	authorGroup := router.Group("/authors")
	{
		authorGroup.POST("/add-author-image", controllers.UploadAuthorImage)
		authorGroup.GET("/:object-key", controllers.RetrieveAuthorImage)
	}

}
