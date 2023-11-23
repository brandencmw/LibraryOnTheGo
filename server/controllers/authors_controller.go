package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"libraryonthego/server/files"
	"libraryonthego/server/services"
	"net/http"
	"strconv"

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

	if err := c.service.AddAuthor(*author); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload author info"})
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
			ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid image option for includeimages provided"))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Option for includeimages was %v, expected boolean value", strImageFlag)})
			return
		}
	}

	id := ctx.Query("id")
	if id == "" {
		c.getAllAuthors(ctx, imageFlag)
	} else {
		id, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid argument for parameter authorID, must be unsigned integer",
			})
			return
		}
		c.getAuthorByID(ctx, uint(id), imageFlag)
	}
}

func (c *AuthorsController) getAllAuthors(ctx *gin.Context, includeImages bool) {
	authors, err := c.service.GetAllAuthors(includeImages)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var image imageResponse
	authorJSON := make([]getAuthorResponse, 0)
	for _, author := range authors {
		fmt.Printf("HEADSHOT: %v\n", author.Headshot)
		if author.Headshot != nil {
			image = imageResponse{
				Name:    author.Headshot.Name,
				Content: base64.StdEncoding.EncodeToString(author.Headshot.Content),
			}
		}
		authorJSON = append(authorJSON, getAuthorResponse{
			ID:        author.ID,
			FirstName: *author.FirstName,
			LastName:  *author.LastName,
			Bio:       *author.Bio,
			Headshot:  image,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"authors": authorJSON})
}

func (c *AuthorsController) getAuthorByID(ctx *gin.Context, ID uint, includeImage bool) {

	author, err := c.service.GetAuthor(ID, includeImage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not retrieve author with ID %v", ID)})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	img := imageResponse{
		Name:    author.Headshot.Name,
		Content: base64.StdEncoding.EncodeToString(author.Headshot.Content),
	}

	resp := getAuthorResponse{
		ID:        author.ID,
		FirstName: *author.FirstName,
		LastName:  *author.LastName,
		Bio:       *author.Bio,
		Headshot:  img,
	}
	ctx.JSON(http.StatusOK, gin.H{"author": resp})
}

func (c *AuthorsController) DeleteAuthor(ctx *gin.Context) {
	strID := ctx.Query("id")
	if strID == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Must have ID in request"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Must have ID in request"})
		return
	}

	ID, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid ID provided"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}

	err = c.service.DeleteAuthor(uint(ID))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to delete: %v", err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *AuthorsController) UpdateAuthor(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Update author info")
}
