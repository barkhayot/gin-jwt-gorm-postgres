package middleware

import (
	"fmt"
	"go-jwt/initializer"
	"go-jwt/model"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// get cookie from request
	tokenString, err := c.Cookie("Autorization")
	if err != nil {
		c.AbortWithStatus(401) // Unautorized
		return
	}

	// decode and validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		//log.Fatal(err)
		c.AbortWithStatus(401) // Unautorized
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check exp time
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(401) // Unautorized
			return
		}
		// check user with token sub
		var user model.User
		initializer.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(401) // Unautorized
			return
		}
		// attach to request
		c.Set("user", user)

		// continue
		c.Next()
	} else {
		c.AbortWithStatus(401) // Unautorized
		return
	}

}
