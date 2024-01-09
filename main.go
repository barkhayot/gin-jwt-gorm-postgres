package main

import (
	"go-jwt/controller"
	"go-jwt/initializer"
	"go-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializer.ConnectDB()
	initializer.LoadEnvVariables()
	initializer.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controller.Signup)
	r.POST("/login", controller.Login)
	r.GET("/secure", middleware.RequireAuth, controller.SecureEndpoint)

	r.Run()
}
