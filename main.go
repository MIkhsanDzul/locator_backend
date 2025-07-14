package main

import (
	"locator-backend/controller"
	"locator-backend/firebase"

	"github.com/gin-gonic/gin"
)

func main() {
	firebase.InitFirebase()

	r := gin.Default()

	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	r.Run(":8080")
}
