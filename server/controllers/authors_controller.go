package controllers

import (
	"fmt"
	"libraryonthego/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthorsController struct {
	service *services.AuthorsService
}

func NewAuthorsController(service *services.AuthorsService) *AuthorsController {
	return &AuthorsController{
		service: service,
	}
}

func (c *AuthorsController) AddAuthor(ctx *gin.Context) {
	var addAuthorRequest addAuthorRequest

	if err := ctx.ShouldBind(&addAuthorRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid request format: %v", err.Error()))
		return
	}

	authorToAdd := services.AddAuthorInfo{
		HeadshotFile: addAuthorRequest.Headshot,
		FirstName:    addAuthorRequest.FirstName,
		LastName:     addAuthorRequest.LastName,
		Bio:          addAuthorRequest.Bio,
	}

	if err := c.service.AddAuthor(authorToAdd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload author info"})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, addAuthorRequest)
}

func (c *AuthorsController) GetAuthor(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Get single author")
}

func (c *AuthorsController) GetAllAuthors(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Get all authors")
}

func (c *AuthorsController) DeleteAuthor(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Delete an author")
}

func (c *AuthorsController) UpdateAuthor(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Update author info")
}
