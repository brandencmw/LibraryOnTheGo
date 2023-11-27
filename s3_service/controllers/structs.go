package controllers

type uploadImageRequest struct {
	Image     []byte `json:"imageContent" binding:"required"`
	ImageName string `json:"imageName" binding:"required"`
}

type replaceImageRequest struct {
	OriginalImageName string `json:"originalName" binding:"required"`
	NewImageName      string `json:"newName" binding:"required"`
	NewImageContent   []byte `json:"newContent"`
}
