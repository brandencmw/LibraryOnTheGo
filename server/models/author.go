package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	ID           int
	FirstName    string
	LastName     string
	Bio          string
	HeadshotPath string
	CreatedAt    datatypes.Time
}
