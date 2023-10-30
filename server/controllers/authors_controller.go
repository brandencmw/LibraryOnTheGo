package controllers

import (
	"errors"
	"libraryonthego/server/services"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func readFormFile(fileHeader *multipart.FileHeader) (fileBytes []byte, err error) {

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	fileBytes = make([]byte, 0)
	for {
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if n == 0 || err != nil {
			break
		}
		fileBytes = append(fileBytes, buffer[:n]...)
	}
	file.Close()

	return fileBytes, err
}

func AddAuthor(c *gin.Context) {

	var addAuthorRequest addAuthorRequest

	err := c.ShouldBind(&addAuthorRequest)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid request format"))
		return
	}

	headshotContents, err := readFormFile(addAuthorRequest.Headshot)
	authorsService := services.DefaultAuthorsService{}
	err = authorsService.AddAuthor(services.AuthorInfo{
		Headshot:  headshotContents,
		FirstName: addAuthorRequest.FirstName,
		LastName:  addAuthorRequest.LastName,
		Bio:       addAuthorRequest.Bio,
	})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

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
