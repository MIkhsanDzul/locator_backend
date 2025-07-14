package main

import (
	"locator-backend/controller"
	"locator-backend/firebase"

	"github.com/gin-gonic/gin"
)

func main() {
	firebase.InitFirebase()

	r := gin.Default()

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.Register)
			auth.POST("/login", controller.Login)
		}	
		location := api.Group("/location")
		{
			location.POST("/update", controller.UpLocation)
			location.GET("/get", controller.GetLocation)
		}
	} 

	r.Run(":8008")
}
