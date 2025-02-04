package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents the authenticated application user
type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex;not null;size:255"`
	PasswordHash  string `gorm:"not null;size:255"`
	LastLogin     *time.Time
	FailedLogins  int `gorm:"default:0"`
	AccountLocked bool `gorm:"default:false"`
	MFAEnabled    bool `gorm:"default:false"`
	MFASecret     string `gorm:"size:255"`
}

// Server represents available SOCKS5 proxy servers
type Server struct {
	gorm.Model
	Country     string  `gorm:"index;not null;size:2"`
	State       string  `gorm:"index;size:100"`
	City        string  `gorm:"index;size:100"`
	Zipcode     string  `gorm:"index;size:20"`
	Address     string  `gorm:"uniqueIndex;not null;size:255"`
	Latency     float64 `gorm:"default:0"`
	IsActive    bool    `gorm:"default:true"`
	LastChecked time.Time
	SSLVerified bool    `gorm:"default:false"`
	Bandwidth   float64 `gorm:"default:0"`
}

// TokenBlacklist tracks revoked JWT tokens
type TokenBlacklist struct {
	gorm.Model
	Token     string    `gorm:"uniqueIndex;not null;size:512"`
	ExpiresAt time.Time `gorm:"index"`
}

// BeforeSave automatically hashes passwords before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	if tx.Statement.Changed("PasswordHash") {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(u.PasswordHash), 
			bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashedPassword)
	}
	return nil
}

// InitSecureDB establishes a secure database connection
func InitSecureDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=verify-full",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: os.Getenv("DB_HOST"),
	}

	dbConfig := postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
		TLS:                  tlsConfig,
	}

	var err error
	db, err = gorm.Open(postgres.New(dbConfig), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect to secure database")
	}

	// Auto-migrate with constraints
	err = db.AutoMigrate(
		&User{},
		&Server{},
		&TokenBlacklist{},
	)
	if err != nil {
		panic("failed to migrate database models")
	}

	// Create partial indexes
	db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower 
		ON users(LOWER(username))
	`)

	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_servers_active 
		ON servers(id) 
		WHERE is_active = true
	`)
}