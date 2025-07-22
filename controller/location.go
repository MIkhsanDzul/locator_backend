package controller

import (
	"context"
	"locator-backend/config"
	"locator-backend/firebase"
	"locator-backend/model"
	"locator-backend/utils"
	"log"
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
	var users []model.User
	if err := config.DB.Table("users").Order("username ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func SaveLocation(c *gin.Context) {
	var loc model.Location
	if err := c.ShouldBindJSON(&loc); err != nil || loc.Username == "" || loc.Latitude == 0 || loc.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request or missing fields"})
		return
	}

	// Cek apakah user ada di Firestore
	docRef := firebase.FirestoreClient.Collection("locations").Doc(loc.Username)
	doc, err := docRef.Get(context.Background())
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	loc.Triggered = false

	// Simpan timestamp ke PostgreSQL
	timestamp := time.Now().UnixMilli()
	timestampFormatted := utils.FormatTimestamp(timestamp)

	err = config.SaveLocationToPostgres(loc.Username, timestamp, timestampFormatted, loc.Latitude, loc.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save location to database"})
		return
	}

	// Ambil data sebelumnya dari Firestore
	existing := doc.Data()
	isRealtime, ok := existing["is_realtime"].(bool)
	if !ok {
		isRealtime = false // default kalau field belum ada
	}

	// Simpan ke Firestore
	firestoreData := map[string]interface{}{
		"latitude":  loc.Latitude,
		"longitude": loc.Longitude,
		"is_realtime": isRealtime,
		"triggered": loc.Triggered,
	}
	// time.Sleep(3 * time.Second)
	_, err = docRef.Set(context.Background(), firestoreData, firestore.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save location to Firestore"})
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
		location.Username = doc.Ref.ID
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
			location.Username = doc.Ref.ID
			locations = append(locations, location)
		}
	}

	return locations, nil
}

func GetLocations(c *gin.Context) {
	ctx := context.Background()
	username := c.Query("username")

	// Ambil data lokasi dari Firestore/PostgreSQL
	locations, err := FetchLocationData(ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch locations"})
		panic(err)
	}

	if username != "" && len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if username == "" && len(locations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no locations found"})
		return
	}

	_, err = firebase.FirestoreClient.
		Collection("locations").
		Doc(username).
		Set(context.Background(), map[string]interface{}{
			"triggered": true,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update Firestore triggered field"})
		return
	}

	// Berhasil
	c.JSON(http.StatusOK, gin.H{"locations": locations})
}


func Realtime(c *gin.Context) {
	username := c.Query("username")
	isRealtime := c.Query("isrealtime") == "true"

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	docRef := firebase.FirestoreClient.Collection("locations").Doc(username)

	// Cek keberadaan dokumen
	docSnap, err := docRef.Get(context.Background())
	if err != nil || !docSnap.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Update is_realtime
	if _, err := docRef.Set(context.Background(), map[string]interface{}{
		"is_realtime": isRealtime,
	}, firestore.MergeAll); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update is_realtime"})
		return
	}

	if isRealtime {
		go func() {
			for {
				// Ambil status terbaru
				snap, err := docRef.Get(context.Background())
				if err != nil || !snap.Exists() {
					log.Println("Failed to fetch latest Firestore document")
					return
				}

				currentRealtime, _ := snap.Data()["is_realtime"].(bool)
				if !currentRealtime {
					// Matikan triggered jika realtime false
					_, _ = docRef.Set(context.Background(), map[string]interface{}{
						"triggered": false,
					}, firestore.MergeAll)
					return
				}

				// Set triggered = true
				_, _ = docRef.Set(context.Background(), map[string]interface{}{
					"triggered": true,
				}, firestore.MergeAll)

				time.Sleep(5 * time.Second)
			}
		}()
	} else {
		// Matikan triggered jika non-realtime
		if _, err := docRef.Set(context.Background(), map[string]interface{}{
			"triggered": false,
		}, firestore.MergeAll); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update triggered to false"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "is_realtime updated",
		"is_realtime": isRealtime,
	})
}

