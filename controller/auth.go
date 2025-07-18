package controller

import (
	"context"
	"locator-backend/firebase"
	"locator-backend/model"
	"locator-backend/utils"

	// "locator-backend/firebase"
	"locator-backend/config"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validasi panjang username
	if len(user.Username) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be at most 16 characters"})
		return
	}

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	// Push baru
	err := config.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func Login(c *gin.Context) {
	var user model.UserLogin
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var dbUser model.User
	err := config.DB.Where("username = ?", user.Username).First(&dbUser).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	if !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// Cek dan update/insert lokasi
	err = config.SaveLocationToPostgres(user.Username, 0, "", 0, 0)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save location to database"})
		return
	}

	// Ubah field "triggered" jadi true di Firestore
	_, err = firebase.FirestoreClient.
		Collection("locations").
		Doc(user.Username).
		Set(context.Background(), map[string]interface{}{
			"triggered": true,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update Firestore triggered field"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
