package main

import (
	"github.com/gin-gonic/gin"
	"go-crud/controllers"
	"go-crud/initializers"
	"go-crud/middleware"
	"net/http"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api is working",
		})
	})
	postsGroup := r.Group("/posts")
	{
		postsGroup.POST("", middleware.RequireAuth, controllers.PostsCreate)
		postsGroup.GET("", controllers.PostsIndex)
		postsGroup.GET("/:id", controllers.PostsShow)
		postsGroup.PUT("/:id", middleware.RequireAuth, controllers.PostsUpdate)
		postsGroup.DELETE("/:id", middleware.RequireAuth, controllers.PostsDelete)
	}

	usersGroup := r.Group("/auth")
	{
		usersGroup.POST("/signup", controllers.SignUp)
		usersGroup.POST("/login", controllers.Login)
		usersGroup.GET("/validate", middleware.RequireAuth, controllers.Validate)
	}

	r.Run()
}
