package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// SFTP configuration
var (
	SFTPHost     string
	SFTPPort     string
	SFTPUser     string
	SFTPPassword string
	SFTPPath     string

	MongoDBConnectionString string
	FileExpiryDays          int
	BaseURL                 string
	Port                    string
)

// LoadEnv loads environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults or environment variables")
	}

	// SFTP Settings - Using IP address instead of hostname
	SFTPHost = getEnvWithDefault("SFTP_HOST", "duynas.duckdns.org") // Reverted to hostname
	SFTPPort = getEnvWithDefault("SFTP_PORT", "9333")
	SFTPUser = getEnvWithDefault("SFTP_USER", "admin")
	SFTPPassword = getEnvWithDefault("SFTP_PASSWORD", "Doanhduy2003")
	SFTPPath = getEnvWithDefault("SFTP_PATH", "/DuyData/cloud")

	// Ensure SFTP path starts with forward slash and uses Unix-style separators
	SFTPPath = filepath.ToSlash(getEnvWithDefault("SFTP_PATH", "/DuyData/cloud"))
	if !strings.HasPrefix(SFTPPath, "/") {
		SFTPPath = "/" + SFTPPath
	}

	// Other settings
	MongoDBConnectionString = getEnvWithDefault("MONGODB_CONNECTION_STRING",
		"mongodb+srv://anhduyking:doanhduy2003@cluster0cloudfrog.e7b3g.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0cloudfrog")
	BaseURL = getEnvWithDefault("BASE_URL", "http://localhost:8080")

	// Parse file expiry days
	FileExpiryDays = 7 // Default 7 days

	// Add port configuration
	Port = getEnvWithDefault("PORT", "3000")
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
