package authentication

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func LoginUser(c *gin.Context) {

	var requestBody authRequestBody

	//extract user and pass from req
	err := c.ShouldBindJSON(&requestBody)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
	}

	//check user and pass against env
	if requestBody.username != os.Getenv("SERVER_ADMIN") || requestBody.password != os.Getenv("SERVER_PASS") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect credentials provided",
		})
	}

	//issue token for future requests
	authToken := jwt.New(jwt.SigningMethodHS256)
	claims := authToken.Claims.(jwt.MapClaims)
	claims["user"] = requestBody.username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := authToken.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to generate token",
		})
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
