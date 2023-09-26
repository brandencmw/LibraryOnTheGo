package authentication

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateUser(c *gin.Context) {

	var authToken string

	err := c.ShouldBindJSON(authToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Must provide token",
		})
	}

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method.Alg())
		}

		return os.Getenv("JWT_SECRET"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().Unix() > claims["exp"].(int64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("message", "authorized")
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
