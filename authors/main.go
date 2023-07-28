package authors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func GetAuthors(c *gin.Context) {
	var data = author{ID: 1, Name: "Guy Fawkes", Bio: "Tried to burn down some building or something"}
	c.IndentedJSON(http.StatusOK, data)
}

func GetAuthor(c *gin.Context) {
	id := c.Param("id")
	c.String(http.StatusOK, id)
}
