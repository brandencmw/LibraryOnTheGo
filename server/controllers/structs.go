package controllers

import "mime/multipart"

type addAuthorRequest struct {
	Headshot  *multipart.FileHeader `form:"headshot" binding:"required"`
	FirstName string                `form:"firstName" binding:"required"`
	LastName  string                `form:"lastName" binding:"required"`
	Bio       string                `form:"bio" binding:"required"`
}

type getAuthorResponse struct {
	ID        uint          `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Bio       string        `json:"bio"`
	Headshot  imageResponse `json:"headshot"`
}

type imageResponse struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}
