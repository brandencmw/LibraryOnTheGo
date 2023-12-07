package controllers

import "mime/multipart"

type addAuthorRequest struct {
	Headshot  *multipart.FileHeader `form:"headshot" binding:"required"`
	FirstName string                `form:"firstName" binding:"required"`
	LastName  string                `form:"lastName" binding:"required"`
	Bio       string                `form:"bio" binding:"required"`
}

type updateAuthorRequest struct {
	Headshot  *multipart.FileHeader `form:"headshot"`
	FirstName *string               `form:"firstName"`
	LastName  *string               `form:"lastName"`
	Bio       *string               `form:"bio"`
	ID        *string               `form:"id"`
}

type getAuthorResponse struct {
	ID        string `json:"id" binding:"required"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
	Headshot  string `json:"headshotKey"`
}

type addBookRequest struct {
	Title       string               `form:"title" binding:"required"`
	Synopsis    string               `form:"synopsis"`
	PublishDate string               `form:"publishDate" binding:"required"`
	PageCount   int                  `form:"pageCount" binding:"required"`
	Categories  []string             `form:"categories"`
	Authors     []string             `form:"authors"`
	Cover       multipart.FileHeader `form:"cover" binding:"required"`
}
