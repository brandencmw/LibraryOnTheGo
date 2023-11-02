package controllers

type uploadImageRequest struct {
	Image     []byte `json:"image" binding:"required"`
	ImageName string `json:"imageName" binding:"required"`
}
