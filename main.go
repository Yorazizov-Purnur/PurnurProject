package main

import (
	"purnur/pjt/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	var r = gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login",controllers.Login)

	r.Run(":2929")
}
