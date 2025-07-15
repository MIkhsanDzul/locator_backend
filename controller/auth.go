package controller

import (
	"context"
	"locator-backend/model"
	"locator-backend/utils"
	"locator-backend/firebase"
	"net/http"

	// "firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// ⛔ Validasi panjang username
	if len(user.Username) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username must be at most 16 characters"})
		return
	}

	ctx := context.Background()
	docRef := firebase.FirestoreClient.Collection("users").Doc(user.Username)

	// 🔍 Cek username sudah ada belum
	doc, err := docRef.Get(ctx)
	if err == nil && doc.Exists() {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	// ✅ Push user baru
	_, err = docRef.Set(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func Login(c *gin.Context) {
	var loginData model.UserLogin
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := context.Background()
	docRef := firebase.FirestoreClient.Collection("users").Doc(loginData.Username)
	doc, err := docRef.Get(ctx)
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var user model.User
	if err := doc.DataTo(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"username": user.Username,
		"email": user.Email,
	})
		
}
