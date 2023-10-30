package controllers

type uploadAuthorImageRequest struct {
	Headshot []byte `json:"image" binding:"required"`
}
