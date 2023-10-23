package controllers

import (
	"fmt"
	"libraryonthego/server/services"
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

	services.SendAuthorImageToS3(addAuthorRequest.Headshot)
	c.JSON(http.StatusOK, addAuthorRequest)
}

func GetAuthor(c *gin.Context) {
	c.String(http.StatusOK, "Get single author")
}

func GetAllAuthors(c *gin.Context) {
	c.String(http.StatusOK, "Get all authors")
}

func DeleteAuthor(c *gin.Context) {
	c.String(http.StatusOK, "Delete an author")
}

func UpdateAuthor(c *gin.Context) {
	c.String(http.StatusOK, "Update author info")
}
