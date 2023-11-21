package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"s3/services"

	"github.com/gin-gonic/gin"
)

type AuthorsController struct {
	Service *services.LibraryBucketService
}

func NewAuthorsControlller(service *services.LibraryBucketService) *AuthorsController {
	return &AuthorsController{
		Service: service,
	}
}

func (c *AuthorsController) UploadAuthorImage(ctx *gin.Context) {
	var request uploadImageRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request format", "uploaded": false})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Error binding request data: %v", err.Error()))
		return
	}

	err = c.Service.UploadImage(request.ImageName, request.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "uploaded": false})
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to upload author image: %v", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uploaded": true,
	})
}

func (c *AuthorsController) RetrieveAuthorImage(ctx *gin.Context) {

	key := ctx.Query("img-name")
	if key == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("Key required for request"))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Key required for request",
		})
		return
	}

	image, err := c.Service.GetImage(key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"imageContent": image.Content,
		"imageName":    image.Name,
	})
}
