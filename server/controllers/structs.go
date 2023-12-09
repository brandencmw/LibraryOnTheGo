package controllers

import "mime/multipart"

type authorRequest struct {
	Headshot  *multipart.FileHeader `form:"headshot" binding:"required"`
	FirstName *string               `form:"firstName" binding:"required"`
	LastName  *string               `form:"lastName" binding:"required"`
	Bio       *string               `form:"bio" binding:"required"`
	ID        *string               `form:"id"`
}

type authorResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
	Headshot  string `json:"headshotKey"`
}

type bookRequest struct {
	Title       string               `form:"title" binding:"required"`
	Synopsis    string               `form:"synopsis"`
	PublishDate string               `form:"publishDate" binding:"required"`
	PageCount   int                  `form:"pageCount" binding:"required"`
	Categories  []string             `form:"categories"`
	Authors     []string             `form:"authors"`
	Cover       multipart.FileHeader `form:"cover" binding:"required"`
	ID          string               `form:"id"`
}

type bookResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Synopsis    string           `json:"synopsis"`
	PublishDate string           `json:"publishDate"`
	PageCount   int              `json:"pageCount"`
	Categories  []string         `json:"categories"`
	Authors     []authorResponse `json:"authors"`
	Cover       string           `json:"coverKey"`
}
