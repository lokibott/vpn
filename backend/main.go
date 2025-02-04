package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	db        *gorm.DB
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

func main() {
	// Load environment variables
	LoadEnv()

	// Initialize database with TLS
	InitSecureDB()
	
	// Create Gin router with security middleware
	r := gin.Default()
	r.Use(SecurityHeadersMiddleware())
	r.Use(limiter.RateLimit())
	
	// Routes
	r.POST("/login", RateLimitedLoginHandler)
	r.GET("/servers", AuthMiddleware(), GetServersHandler)
	
	// Start HTTPS server
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	srv := &http.Server{
		Addr:         ":443",
		Handler:      r,
		TLSConfig:    cfg,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Redirect HTTP to HTTPS
	go func() {
		log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS)))
	}()

	log.Fatal(srv.ListenAndServeTLS(
		os.Getenv("TLS_CERT"),
		os.Getenv("TLS_KEY"),
	))
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
}