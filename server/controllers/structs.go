package controllers

import "mime/multipart"

type addAuthorRequest struct {
	Headshot  *multipart.FileHeader `form:"headshot" binding:"required"`
	FirstName string                `form:"firstName" binding:"required"`
	LastName  string                `form:"lastName" binding:"required"`
	Bio       string                `form:"bio" binding:"required"`
}
