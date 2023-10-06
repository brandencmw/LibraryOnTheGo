package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAuthor(c *gin.Context) {

	var addAuthorRequest addAuthorRequest

	c.Request.ParseMultipartForm(1000000)
	fileHeader, _ := c.FormFile("headshot")
	file, _ := fileHeader.Open()

	fileBytes := make([]byte, 0)
	for {
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break
		}
		fileBytes = append(fileBytes, buffer[:n]...)
	}
	file.Close()

	addAuthorRequest.Headshot = fileBytes
	addAuthorRequest.FirstName = c.PostForm("firstName")
	addAuthorRequest.LastName = c.PostForm("lastName")
	addAuthorRequest.Bio = c.PostForm("bio")

	fmt.Printf("Author Data: %v\n", addAuthorRequest)
	c.JSON(http.StatusOK, addAuthorRequest)
}
