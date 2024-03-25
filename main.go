package main

import (
	"purnur/pjt/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	var r = gin.Default()

	r.Use(Cors)

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/allposts", controllers.AllPosts)
	r.GET("/mypost", controllers.MyPost)
	r.POST("/createpost", controllers.CreateFunc)
	r.POST("/like",controllers.Like)

	r.Run(":2929")
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.75:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}

