package controller

import (
	"context"
	"locator-backend/firebase"
	"locator-backend/model"
	"locator-backend/utils"
	"strconv"
	"time"

	// "locator-backend/utils"
	// "log"
	"net/http"

	// "firebase.google.com/go/v4/db"
	"cloud.google.com/go/firestore"
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

	// Cek username ada atau tidak
	docRef := firebase.FirestoreClient.Collection("users").Doc(loc.Username)
	doc, err := docRef.Get(context.Background())
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	loc.Timestamp = time.Now().UnixMilli()
	loc.TimestampFormatted = utils.FormatTimestamp(loc.Timestamp)
	loc.Triggered = false

	//Push baru 
	_, err = firebase.FirestoreClient.Collection("locations").Doc(loc.Username).Set(context.Background(), loc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location saved"})
}

func FetchLocationData(ctx context.Context, username string) ([]model.Location, error) {
	if username != "" {
		docRef := firebase.FirestoreClient.Collection("locations").Doc(username)
		doc, err := docRef.Get(ctx)
		if err != nil || !doc.Exists() {
			return nil, err
		}

		var location model.Location
		if err := doc.DataTo(&location); err != nil {
			return nil, err
		}
		location.Triggered = true
		return []model.Location{location}, nil
	}

	// fetch semua
	docs, err := firebase.FirestoreClient.Collection("locations").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	locations := make([]model.Location, 0)
	for _, doc := range docs {
		var location model.Location
		if err := doc.DataTo(&location); err == nil {
			locations = append(locations, location)
		}
	}

	return locations, nil
}

func GetLocations(c *gin.Context) {
	username := c.Query("username")
	ctx := context.Background()

	locations, err := FetchLocationData(ctx, username)
	if err != nil || len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found or no location data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"locations": locations})
}

func Realtime(c *gin.Context) {
	ctx := context.Background()
	username := c.Query("username")
	isRealtimeStr := c.DefaultQuery("isrealtime", "true")
	locationsCollection := firebase.FirestoreClient.Collection("locations")

	// Parse isrealtime string ke bool
	isRealtime, err := strconv.ParseBool(isRealtimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid isrealtime value (must be true/false)"})
		return
	}

	//Fungsi update dokumen
	updateRealtime := func(ref *firestore.DocumentRef) error {
		_, err := ref.Update(ctx, []firestore.Update{
			{Path: "IsRealtime", Value: isRealtime},
		})
		return err
	}

	//Update satu user atau semua
	if username != "" {
		docRef := locationsCollection.Doc(username)
		docSnap, err := docRef.Get(ctx)
		if err != nil || !docSnap.Exists() {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if err := updateRealtime(docRef); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update is_realtime"})
			return
		}
	} else {
		docs, err := locationsCollection.Documents(ctx).GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch documents"})
			return
		}
		for _, doc := range docs {
			_ = updateRealtime(doc.Ref)
		}
	}

	//Ambil data lokasi
	locations, err := FetchLocationData(ctx, username)
	if err != nil || len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no location data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"locations": locations,
	})
}


