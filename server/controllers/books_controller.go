package controllers

import (
	"context"
	"errors"
	"fmt"
	"libraryonthego/server/files"
	"libraryonthego/server/services"
	"mime/multipart"
	"net/http"
	"strconv"

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

func requestImageToServiceImage(img multipart.FileHeader) (*services.Image, error) {
	imageContent, err := files.GetMultipartFormContents(&img)
	if err != nil {
		return nil, err
	}
	return &services.Image{
		Content: imageContent,
		Name:    img.Filename,
	}, nil
}

func collectBookOptionsFromRequest(req bookRequest) ([]services.BookOption, error) {
	var options []services.BookOption
	if req.Title != nil {
		options = append(options, services.WithTitle(*req.Title))
	}
	if req.Synopsis != nil {
		options = append(options, services.WithSynopsis(*req.Synopsis))
	}
	if req.PublishDate != nil {
		options = append(options, services.WithPublishDate(*req.PublishDate))
	}
	if req.PageCount != nil {
		options = append(options, services.WithPageCount(*req.PageCount))
	}
	if len(req.Categories) > 0 {
		options = append(options, services.WithCategories(req.Categories...))
	}
	if len(req.Authors) > 0 {
		options = append(options, services.WithAuthors(req.Authors...))
	}
	if req.Cover != nil {
		cover, err := requestImageToServiceImage(*req.Cover)
		if err != nil {
			return nil, err
		}
		options = append(options, services.WithCover(*cover))
	}
	return options, nil
}

func bookRequestToServiceBook(req bookRequest) (*services.Book, error) {
	options, err := collectBookOptionsFromRequest(req)
	if err != nil {
		return nil, err
	}
	if len(options) == 0 {
		return nil, errors.New("Must provide some info in request")
	}
	return services.NewBook(options...)
}

func serviceBookToBookResponse(book *services.BookOutput) bookResponse {
	authors := make([]authorResponse, 0, len(book.Authors))
	for _, author := range book.Authors {
		authors = append(authors, serviceAuthorToAuthorResponse(&author))
	}
	return bookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Synopsis:    book.Synopsis,
		PublishDate: book.PublishDate.Format("2006-01-02"),
		PageCount:   book.PageCount,
		Cover:       book.CoverObjectKey,
		Categories:  book.Categories,
		Authors:     authors,
	}
}

func (c *BooksController) AddBook(ctx *gin.Context) {
	var req bookRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid request format: %v", err.Error()))
		return
	}

	book, err := bookRequestToServiceBook(req)
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
	strImageFlag := ctx.Query("includeimages")
	var imageFlag bool
	var err error
	if strImageFlag == "" {
		imageFlag = true
	} else {
		imageFlag, err = strconv.ParseBool(strImageFlag)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid image option for includeimages provided, received %v", strImageFlag))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Option for includeimages was %v, expected boolean value", strImageFlag)})
			return
		}
	}

	ID := ctx.Query("id")
	if ID == "" {
		c.getAllBooks(ctx, imageFlag)
	} else {
		c.getBookByID(ctx, ID, imageFlag)
	}
}

func (c *BooksController) getAllBooks(ctx *gin.Context, includeImages bool) {
	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	books, err := c.service.GetAllBooks(parent, includeImages)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
			ctx.AbortWithError(http.StatusRequestTimeout, fmt.Errorf("Server timed out: %v", err))
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get books"})
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Could not get books: %v", err))
		}
	}

	resBooks := make([]bookResponse, 0, len(books))
	for _, book := range books {
		resBooks = append(resBooks, serviceBookToBookResponse(book))
	}
	ctx.JSON(http.StatusOK, gin.H{"books": resBooks})
}

func (c *BooksController) getBookByID(ctx *gin.Context, ID string, includeImage bool) {
	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	book, err := c.service.GetBook(parent, ID, includeImage)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
			ctx.AbortWithError(http.StatusRequestTimeout, fmt.Errorf("Server timed out: %v", parent.Err().Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not get book with ID %v", ID)})
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Could not get book with ID %v: %v", ID, err))
		}
	}
	res := serviceBookToBookResponse(book)
	ctx.JSON(http.StatusOK, gin.H{"book": res})
}

func (c *BooksController) DeleteBook(ctx *gin.Context) {

}

func (c *BooksController) UpdateBook(ctx *gin.Context) {

	var req bookRequest
	ctx.ShouldBind(&req)
	if req.ID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Need ID for update request"})
		ctx.AbortWithError(http.StatusBadRequest, errors.New("No ID provided for update request"))
		return
	}

	book, err := bookRequestToServiceBook(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create book from request"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to create book: %v", err.Error()))
		return
	}

	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	err = c.service.UpdateBook(parent, *req.ID, *book)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
			ctx.AbortWithError(http.StatusRequestTimeout, fmt.Errorf("Server timed out: %v", parent.Err().Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update book"})
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to update book: %v", err.Error()))
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
