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

// Firestore
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
	username := c.Query("username")
	ctx := context.Background()

	if username != "" {
		docRef := firebase.FirestoreClient.Collection("locations").Doc(username)
		doc, err := docRef.Get(ctx)
		if err != nil || !doc.Exists() {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var location model.Location
		if err := doc.DataTo(&location); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse location data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"location": location})
		return
	}

	// Get all locations
	docs, err := firebase.FirestoreClient.Collection("locations").Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch locations"})
		return
	}

	locations := make([]model.Location, len(docs))
	for i, doc := range docs {
		var location model.Location
		if err := doc.DataTo(&location); err == nil {
			locations[i] = location
		}
	}

	c.JSON(http.StatusOK, gin.H{"locations": locations})
}
