package controllers

import (
	"context"
	"errors"
	"fmt"
	"libraryonthego/server/files"
	"libraryonthego/server/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const requestTimeout = time.Second * 2

type AuthorsController struct {
	service *services.AuthorsService
}

func NewAuthorsController(service *services.AuthorsService) *AuthorsController {
	return &AuthorsController{
		service: service,
	}
}

func (c *AuthorsController) AddAuthor(ctx *gin.Context) {
	var req addAuthorRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid request format: %v", err.Error()))
		return
	}

	imageContent, err := files.GetMultipartFormContents(req.Headshot)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Could not read file contents"})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	author, err := services.NewAuthor(
		services.WithFirstName(req.FirstName),
		services.WithLastName(req.LastName),
		services.WithBio(req.Bio),
		services.WithImage(&services.Image{Name: req.Headshot.Filename, Content: imageContent}),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create author"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to create author: %v", err))
		return
	}

	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	if err = c.service.AddAuthor(parent, *author); err != nil {
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

func (c *AuthorsController) GetAuthor(ctx *gin.Context) {

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
		c.getAllAuthors(ctx, imageFlag)
	} else {
		c.getAuthorByID(ctx, ID, imageFlag)
	}
}

func (c *AuthorsController) getAllAuthors(ctx *gin.Context, includeImages bool) {
	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	authors, err := c.service.GetAllAuthors(parent, includeImages)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authors"})
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	authorJSON := make([]getAuthorResponse, 0)
	for _, author := range authors {
		authorJSON = append(authorJSON, getAuthorResponse{
			ID:        author.ID,
			FirstName: author.FirstName,
			LastName:  author.LastName,
			Bio:       author.Bio,
			Headshot:  author.HeadshotObjectKey,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"authors": authorJSON})
}

func (c *AuthorsController) getAuthorByID(ctx *gin.Context, ID string, includeImage bool) {

	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	author, err := c.service.GetAuthor(parent, ID, includeImage)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not retrieve author with ID %v", ID)})
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp := getAuthorResponse{
		ID:        author.ID,
		FirstName: author.FirstName,
		LastName:  author.LastName,
		Bio:       author.Bio,
		Headshot:  author.HeadshotObjectKey,
	}
	ctx.JSON(http.StatusOK, gin.H{"author": resp})
}

func (c *AuthorsController) DeleteAuthor(ctx *gin.Context) {
	ID := ctx.Query("id")
	if ID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Must have ID in request"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Must have ID in request"})
		return
	}

	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	err := c.service.DeleteAuthor(parent, ID)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to delete: %v", err.Error()))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func collectUpdateAuthorOptions(req updateAuthorRequest) ([]services.AuthorOption, error) {
	options := make([]services.AuthorOption, 0)
	if req.FirstName != nil {
		options = append(options, services.WithFirstName(*req.FirstName))
	}
	if req.LastName != nil {
		options = append(options, services.WithLastName(*req.LastName))
	}
	if req.Bio != nil {
		options = append(options, services.WithBio(*req.Bio))
	}
	if req.Headshot != nil {
		imageContent, err := files.GetMultipartFormContents(req.Headshot)
		if err != nil {
			return nil, err
		}
		options = append(options, services.WithImage(&services.Image{Name: req.Headshot.Filename, Content: imageContent}))
	}
	return options, nil
}

func (c *AuthorsController) UpdateAuthor(ctx *gin.Context) {
	var req updateAuthorRequest
	ctx.ShouldBind(&req)
	if req.ID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Must provide ID of author to update"})
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Must provide ID of author to update"))
		return
	}

	options, err := collectUpdateAuthorOptions(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request"})
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	if len(options) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Must have at least one field to update"})
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Must have at least one field to update"))
		return
	}

	author, err := services.NewAuthor(options...)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create author"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to create author: %v", err))
		return
	}

	parent, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	err = c.service.UpdateAuthor(parent, *req.ID, *author)
	if err != nil {
		if parent.Err() != nil {
			ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Server took too long to respond"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to update author: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, req)
}
