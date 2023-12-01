package controllers

import (
	"context"
	"libraryonthego/googlebooks/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BooksController struct {
	service *services.BooksService
}

type SearchBooksResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	PublishDate string   `json:"publishYear"`
	Description string   `json:"description"`
	PageCount   uint     `json:"pageCount"`
	Categories  []string `json:"categories"`
}

func NewBooksController(service *services.BooksService) *BooksController {
	return &BooksController{
		service: service,
	}
}

func (c *BooksController) GetBookInfo(ctx *gin.Context) {
	title := ctx.Query("title")
	author := ctx.Query("author")

	if title == "" && author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Need title or author to make request"})
	}

	parent, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	books, err := c.service.SearchBooks(parent, title, author)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	bookJSON := make([]SearchBooksResponse, 0, len(books))
	for _, book := range books {
		bookJSON = append(bookJSON, SearchBooksResponse{
			ID:          book.ID,
			Title:       book.Title,
			Authors:     book.Authors,
			PublishDate: book.PublishDate,
			Description: book.Description,
			PageCount:   book.PageCount,
			Categories:  book.Categories,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"books": bookJSON})
}
