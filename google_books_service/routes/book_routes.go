package routes

import (
	"libraryonthego/googlebooks/controllers"

	"github.com/gin-gonic/gin"
)

func AttachBookRotues(router *gin.Engine, controller *controllers.BooksController) {
	router.GET("/book", controller.GetBookInfo)
}