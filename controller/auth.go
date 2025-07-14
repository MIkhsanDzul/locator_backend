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

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	ctx := context.Background()
	client, _ := firebase.App.DatabaseWithURL(ctx, "https://locator-dccf6-default-rtdb.asia-southeast1.firebasedatabase.app/")
	ref := client.NewRef("users")

	_, err := ref.Push(ctx, user)
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
		if u.Email == loginData.Email && utils.CheckPasswordHash(loginData.Password, u.Password) {
			c.JSON(http.StatusOK, gin.H{"message": "login successful", "name": u.Username})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}
