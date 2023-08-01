package main

import (
	"libraryonthego/server/config"
	"libraryonthego/server/models"
)

func init() {
	config.LoadEnv()
	config.DBInit()
}

func main() {
	config.DB.AutoMigrate(&models.Author{})
}
