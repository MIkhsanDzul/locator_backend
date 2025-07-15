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
	client, err := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "firebase init failed"})
		return
	}

	// 🔍 Cek username sudah ada belum
	usersRef := client.NewRef("users")
	var users map[string]map[string]interface{}
	if err := usersRef.Get(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing users"})
		return
	}

	for _, u := range users {
		if uname, ok := u["username"].(string); ok && uname == user.Username {
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
			return
		}
	}

	// 🔐 Hash password
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	// ✅ Push user baru
	_, err = usersRef.Push(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func Login(c *gin.Context) {
	var loginData model.User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ctx := context.Background()
	client, _ := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	ref := client.NewRef("users")

	var users map[string]model.User
	if err := ref.Get(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	for _, u := range users {
		if u.Username == loginData.Username && utils.CheckPasswordHash(loginData.Password, u.Password) {
			c.JSON(http.StatusOK, gin.H{"message": "login successful", "name": u.Username})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}
