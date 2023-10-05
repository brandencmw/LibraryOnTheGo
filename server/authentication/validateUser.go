package authentication

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateUser(c *gin.Context) {

	var httpStatus int

	authorized, exists := c.Get("authorized")

	fmt.Printf("Authenticated: %v\n", authorized)
	fmt.Printf("Exists: %v\n", exists)

	if authorized == "true" {
		httpStatus = http.StatusOK
	} else {
		httpStatus = http.StatusUnauthorized
	}

	c.JSON(httpStatus, gin.H{
		"authenticated": authorized,
		"exists":        exists,
	})
}
