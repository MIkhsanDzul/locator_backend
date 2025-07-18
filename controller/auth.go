package controller

import (
	"locator-backend/model"
	"locator-backend/utils"
	// "locator-backend/firebase"
	"net/http"
	"locator-backend/config"

	"github.com/gin-gonic/gin"
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

	location, err := FetchLocationData(c.Request.Context(), user.Username)
	if len(location) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found or no location data"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch location"})
		panic(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "location": location})
}

