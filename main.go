package main

import (
	"locator-backend/controller"
	"locator-backend/firebase"
	"locator-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	firebase.InitFirestore()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controller.Register)
			auth.POST("/login", controller.Login)
			auth.GET("/users", controller.GetUsers)
		}	
		location := api.Group("/location")
		{
			location.POST("/update", controller.SaveLocation)
			location.GET("/get", controller.GetLocations)
			location.GET("/realtime", controller.Realtime)
		}
	} 

	r.Run(":8008")
}
