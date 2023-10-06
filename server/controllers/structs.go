package controllers

type addAuthorRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Bio       string `json:"bio" binding:"required"`
	Headshot  []byte `json:"headshot" binding:"required"`
}
