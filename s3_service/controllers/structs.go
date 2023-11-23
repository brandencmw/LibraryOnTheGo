package controllers

type uploadImageRequest struct {
	Image     []byte `json:"imageContent" binding:"required"`
	ImageName string `json:"imageName" binding:"required"`
}
