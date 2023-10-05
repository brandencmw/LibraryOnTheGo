package authentication

import (
	"fmt"
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

	fmt.Println(c.Request.Body)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid request format. Received: %v", requestBody))
		return
	}

	//check user and pass against env
	if requestBody.Username != os.Getenv("SERVER_ADMIN") || requestBody.Password != os.Getenv("SERVER_PASS") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect credentials provided",
		})
		return
	}

	//issue token for future requests
	authToken := jwt.New(jwt.SigningMethodHS256)
	claims := authToken.Claims.(jwt.MapClaims)
	claims["user"] = requestBody.Username
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	secret := []byte(os.Getenv(("JWT_SECRET")))
	tokenString, err := authToken.SignedString(secret)
	if err != nil {
		c.AbortWithError(http.StatusInsufficientStorage, fmt.Errorf("Unable to generate token with %v", secret))
		c.AbortWithError(http.StatusBadGateway, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
