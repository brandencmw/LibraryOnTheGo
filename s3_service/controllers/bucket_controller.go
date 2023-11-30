package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"s3/services"

	"github.com/gin-gonic/gin"
)

type BucketController struct {
	Bucket  string
	Service *services.BucketService
}

func NewBucketController(bucket string, service *services.BucketService) *BucketController {
	return &BucketController{
		Bucket:  bucket,
		Service: service,
	}
}

func (c *BucketController) UploadImage(ctx *gin.Context) {
	var request uploadImageRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect request format", "uploaded": false})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Error binding request data: %v", err.Error()))
		return
	}

	err = c.Service.UploadImage(ctx, c.Bucket, request.ImageName, request.Image)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "uploaded": false})
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to upload author image: %v", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"uploaded": true,
	})
}

func (c *BucketController) DeleteImage(ctx *gin.Context) {
	imageName := ctx.Query("img-name")
	if imageName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request requires img-name param"})
		ctx.AbortWithError(http.StatusBadRequest, errors.New("img-name param not provided"))
		return
	}
	err := c.Service.DeleteImage(ctx, c.Bucket, imageName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *BucketController) RetrieveObjectKey(ctx *gin.Context) {
	img := ctx.Query("img")
	if img == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Must have img parameter to retrieve key"})
		ctx.AbortWithError(http.StatusBadRequest, errors.New("img parameter not specified"))
		return
	}

	key, err := c.Service.GetObjectKey(ctx, c.Bucket, img)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Object with name %v not found", img)})
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"key": key})
}

func (c *BucketController) ReplaceImage(ctx *gin.Context) {
	var request replaceImageRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = c.Service.ReplaceImage(ctx, c.Bucket, request.OriginalImageName, request.NewImageName, request.NewImageContent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image"})
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to update image: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
