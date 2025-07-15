package controller

import (
	"context"
	"locator-backend/firebase"
	"locator-backend/model"
	"locator-backend/utils"
	"time"

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
	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
		return
	}

	// 🔍 Ambil semua user
	usersRef := client.NewRef("users")
	var users map[string]map[string]interface{}
	if err := usersRef.Get(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	// 🔎 Cari apakah username-nya ada
	found := false
	for _, user := range users {
		if uname, ok := user["username"].(string); ok && uname == location.Username {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// ✅ Simpan lokasi
	location.Timestamp = time.Now().UnixMilli()
	location.TimestampFormatted = utils.FormatTimestamp(location.Timestamp) //Abdi format

	ref := client.NewRef("locations/" + location.Username)
	if err := ref.Set(ctx, location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
}

func GetLocation(c *gin.Context) {
	ctx := context.Background()
	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
		return
	}

	// Ambil query
	username := c.Query("username")

	// Kalau ada username, ambil satu lokasi
	if username != "" {
		// Cek apakah username ada di /users
		usersRef := client.NewRef("users")
		var users map[string]map[string]interface{}
		if err := usersRef.Get(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
			return
		}

		found := false
		for _, user := range users {
			if uname, ok := user["username"].(string); ok && uname == username {
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Ambil lokasi spesifik
		ref := client.NewRef("locations/" + username)
		var location model.Location
		if err := ref.Get(ctx, &location); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch location"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"location": location})
		return
	}

	// Kalau gak ada username: ambil semua lokasi (limit 10)
	locRef := client.NewRef("locations")
	var allLocations map[string]model.Location
	if err := locRef.Get(ctx, &allLocations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch all locations"})
		return
	}

	// Limit 10
	result := make([]model.Location, 0, 10)
	count := 0
	for _, loc := range allLocations {
		result = append(result, loc)
		count++
		if count >= 10 {
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"locations": result})
}

func GetUsers(c *gin.Context) {
	ctx := context.Background()
	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
		return
	}

	usersRef := client.NewRef("users")
	var users map[string]map[string]interface{}
	if err := usersRef.Get(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}