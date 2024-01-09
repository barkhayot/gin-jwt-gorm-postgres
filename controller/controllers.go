package controller

import (
	"go-jwt/initializer"
	"go-jwt/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	// Get user data from req
	var body struct { // expected data
		Email    string
		Password string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash password"})
		return
	}

	// Record data
	user := model.User{Email: body.Email, Password: string(hash)}
	result := initializer.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Failed to create user"})
		return
	}

	// Return 200
	c.JSON(200, gin.H{"message": "User has been created"})

}

func Login(c *gin.Context) {
	// get data from req
	var body struct { // expected data
		Email    string
		Password string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// look up user data
	var user model.User
	initializer.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(404, gin.H{"error": "User with that email is not found"})
		return
	}

	// compare hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid email or password"})
		return
	}

	// generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to generate token"})
		return
	}

	// return token
	//c.JSON(200, gin.H{"token": tokenString})

	// set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Autorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{})

}

func SecureEndpoint(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(200, gin.H{"user": user})
}

func CheckData(c *gin.Context) {
	var users []model.User
	initializer.DB.Find(&users)

	// return users
	c.JSON(200, gin.H{"users": users})

}
