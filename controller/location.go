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

//Firestore
func GetUsers(c *gin.Context) {
	username := c.Query("username")
	ctx := context.Background()

	if username != "" {
		docRef := firebase.FirestoreClient.Collection("users").Doc(username)
		doc, err := docRef.Get(ctx)
		if err != nil || !doc.Exists() {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var user model.User
		if err := doc.DataTo(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
		return
	}

	// Get all users
	docs, err := firebase.FirestoreClient.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	users := make([]model.User, len(docs))
	for i, doc := range docs {
		var user model.User
		if err := doc.DataTo(&user); err == nil {
			users[i] = user
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func SaveLocation(c *gin.Context) {
	var loc model.Location
	if err := c.ShouldBindJSON(&loc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// ⛔ Validasi panjang username
	if len(loc.Username) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be at most 16 characters"})
		return
	}

	// 🔍 Cek username sudah ada belum
	docRef := firebase.FirestoreClient.Collection("users").Doc(loc.Username)
	doc, err := docRef.Get(context.Background())
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	loc.Timestamp = time.Now().UnixMilli()
	loc.TimestampFormatted = utils.FormatTimestamp(loc.Timestamp)

	// ✅ Push user baru
	_, err = firebase.FirestoreClient.Collection("locations").Doc(loc.Username).Set(context.Background(), loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location saved"})
}

func GetLocations(c *gin.Context) {
	locations, err := firebase.FirestoreClient.Collection("locations").Documents(context.Background()).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch locations"})
		return
	}

	result := make([]model.Location, len(locations))
	for i, doc := range locations {
		var loc model.Location
		doc.DataTo(&loc)
		result[i] = loc
	}

	c.JSON(http.StatusOK, gin.H{"locations": result})
}

// func UpLocation(c *gin.Context) {
// 	var location model.Location
// 	if err := c.ShouldBindJSON(&location); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	ctx := context.Background()
// 	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
// 		return
// 	}

// 	// 🔍 Ambil semua user
// 	usersRef := client.NewRef("users")
// 	var users map[string]map[string]interface{}
// 	if err := usersRef.Get(ctx, &users); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
// 		return
// 	}

// 	// 🔎 Cari apakah username-nya ada
// 	found := false
// 	for _, user := range users {
// 		if uname, ok := user["username"].(string); ok && uname == location.Username {
// 			found = true
// 			break
// 		}
// 	}

// 	if !found {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 		return
// 	}

// 	// ✅ Simpan lokasi
// 	location.Timestamp = time.Now().UnixMilli()
// 	location.TimestampFormatted = utils.FormatTimestamp(location.Timestamp) //Abdi format

// 	ref := client.NewRef("locations/" + location.Username)
// 	if err := ref.Set(ctx, location); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update location"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
// }

// func GetLocation(c *gin.Context) {
// 	ctx := context.Background()
// 	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
// 		return
// 	}

// 	// Ambil query
// 	username := c.Query("username")

// 	// Kalau ada username, ambil satu lokasi
// 	if username != "" {
// 		// Cek apakah username ada di /users
// 		usersRef := client.NewRef("users")
// 		var users map[string]map[string]interface{}
// 		if err := usersRef.Get(ctx, &users); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
// 			return
// 		}

// 		found := false
// 		for _, user := range users {
// 			if uname, ok := user["username"].(string); ok && uname == username {
// 				found = true
// 				break
// 			}
// 		}

// 		if !found {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
// 			return
// 		}

// 		// Ambil lokasi spesifik
// 		ref := client.NewRef("locations/" + username)
// 		var location model.Location
// 		if err := ref.Get(ctx, &location); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch location"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"location": location})
// 		return
// 	}

// 	// Kalau gak ada username: ambil semua lokasi (limit 10)
// 	locRef := client.NewRef("locations")
// 	var allLocations map[string]model.Location
// 	if err := locRef.Get(ctx, &allLocations); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch all locations"})
// 		return
// 	}

// 	// Limit 10
// 	result := make([]model.Location, 0, 10)
// 	count := 0
// 	for _, loc := range allLocations {
// 		result = append(result, loc)
// 		count++
// 		if count >= 10 {
// 			break
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"locations": result})
// }

// func GetUsers(c *gin.Context) {
// 	ctx := context.Background()
// 	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
// 		return
// 	}

// 	usersRef := client.NewRef("users")
// 	var users map[string]map[string]interface{}
// 	if err := usersRef.Get(ctx, &users); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"users": users})
// }