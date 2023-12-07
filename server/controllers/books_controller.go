package controllers

import (
	"context"
	"fmt"
	"libraryonthego/server/files"
	"libraryonthego/server/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BooksController struct {
	service *services.BooksService
}

func NewBooksController(service *services.BooksService) *BooksController {
	return &BooksController{
		service: service,
	}
}

func (c *BooksController) AddBook(ctx *gin.Context) {
	var req addBookRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid request format: %v", err.Error()))
		return
	}

	imageContent, err := files.GetMultipartFormContents(&req.Cover)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not read file contents"})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	book, err := services.NewBook(
		services.WithTitle(req.Title),
		services.WithSynopsis(req.Synopsis),
		services.WithAuthors(req.Authors...),
		services.WithPageCount(req.PageCount),
		services.WithPublishDate(req.PublishDate),
		services.WithCategories(req.Categories...),
		services.WithCover(services.Image{Content: imageContent, Name: req.Cover.Filename}),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create book"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to create book: %v", err))
		return
	}

	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	if err = c.service.AddBook(parent, book); err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Request took too long to respond"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload author info"})
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, req)
}

func (c *BooksController) GetBook(ctx *gin.Context) {

}

func (c *BooksController) DeleteBook(ctx *gin.Context) {

}

func (c *BooksController) UpdateBook(ctx *gin.Context) {

}
