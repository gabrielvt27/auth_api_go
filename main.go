package main

import (
	"api-go/controllers"
	"api-go/initializers"
	"api-go/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.GET("/", controllers.HelloApi)
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequireAuth, controllers.Validate)
	r.POST("/edit-email", middlewares.RequireAuth, controllers.EditEmail)
	r.POST("/edit-password", middlewares.RequireAuth, controllers.EditPassword)

	r.Run()
}
