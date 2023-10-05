package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(c *gin.Context) {
	authToken, err := c.Cookie("Authorization")
	fmt.Printf("Cookie:%v\n", authToken)
	
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Cookie:%v", authToken))
		return
	}

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Method)
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	fmt.Printf("Valid: %v\n", token.Valid)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusBadGateway)
		}

		c.Set("authorized", "true")
		c.Next()
	} else {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Wrong"))
	}
}
