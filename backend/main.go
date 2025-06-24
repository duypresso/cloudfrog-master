package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"f/cloudfrog/backend/config"
	"f/cloudfrog/backend/controllers"
	"f/cloudfrog/backend/database"
	"f/cloudfrog/backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// runCleanup is a wrapper function to call AutoCleanupFiles without returning its value
func runCleanup() {
	count := controllers.AutoCleanupFiles()
	log.Printf("Auto cleanup completed: %d expired files removed", count)
}

func main() {
	// Load configuration
	config.LoadEnv()

	// Connect to database
	database.Connect()
	defer database.Disconnect()

	// Initialize router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting middleware
	r.Use(middleware.RateLimiter())

	// Routes
	r.POST("/upload", controllers.UploadFile)
	r.GET("/download/:shortcode", controllers.DownloadFile)
	r.DELETE("/cleanup", middleware.AdminAuth(), controllers.CleanupFiles)

	// Setup cron job for automatic file cleanup
	c := cron.New()
	c.AddFunc("0 0 * * *", runCleanup) // Run daily at midnight
	c.Start()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		c.Stop()
		database.Disconnect()
		os.Exit(0)
	}()

	// Get port from environment variable with fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
