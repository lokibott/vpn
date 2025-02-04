package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Address string `json:"address"`
}

func LoginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Verify credentials
	var dbUser User
	if result := db.Where("username = ?", user.Username).First(&dbUser); result.Error != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, _ := GenerateJWT(dbUser)
	c.JSON(200, gin.H{"token": token})
}

func GetServersHandler(c *gin.Context) {
	// Implementation to fetch servers from vendor API
	var servers []Server
	// Add filtering logic based on query params
	c.JSON(200, servers)
}