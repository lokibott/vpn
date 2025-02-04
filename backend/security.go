package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
)

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return secure.New(secure.Config{
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})
}

// GenerateSecureToken creates a cryptographically secure JWT secret
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// HashPassword uses bcrypt with cost 14
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// LoadEnv handles environment configuration
func LoadEnv() {
	// In production, use actual environment variables
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading .env file")
		}
	}
	
	// Ensure JWT_SECRET exists
	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET environment variable not set")
	}
}