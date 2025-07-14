package controller

import (
	"context"
	"locator-backend/firebase"
	"locator-backend/model"
	// "locator-backend/utils"
	// "log"
	"net/http"

	// "firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

func UpLocation(c *gin.Context) {
	var location model.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := context.Background()
	client, _ := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	ref := client.NewRef("locations")

	_, err := ref.Push(ctx, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update location"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
}

func GetLocation(c *gin.Context) {
	ctx := context.Background()
	client, _ := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	ref := client.NewRef("locations")

	var locations []model.Location
	if err := ref.Get(ctx, &locations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch locations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"locations": locations})
}