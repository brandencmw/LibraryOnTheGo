package main

import (
	"libraryonthego/server/config"
	"libraryonthego/server/models"
)

func init() {
	config.DBInit()
}

func main() {
	config.DB.AutoMigrate(&models.Author{})
}
