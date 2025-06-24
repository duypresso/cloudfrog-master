package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"f/cloudfrog/backend/config"
	"f/cloudfrog/backend/database"
	"f/cloudfrog/backend/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadFile handles file upload requests
func UploadFile(c *gin.Context) {
	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		log.Printf("Error creating temp file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy file data to temp file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Printf("Error copying file data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate unique filename and shortcode
	shortCode, err := utils.GenerateShortCode()
	if err != nil {
		log.Printf("Error generating shortcode: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate shortcode"})
		return
	}

	// Create unique filename (using shortcode to avoid collisions)
	fileExt := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", shortCode, fileExt)
	remotePath := filepath.Join(config.SFTPPath, fileName)

	// Upload to SFTP with detailed error logging
	err = utils.UploadToSFTP(tempFile, remotePath)
	if err != nil {
		log.Printf("SFTP upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// Calculate expiry date
	expiryDate := time.Now().AddDate(0, 0, config.FileExpiryDays)

	// Save to database
	fileRecord := database.FileRecord{
		ID:           primitive.NewObjectID(),
		FileName:     fileName,
		OriginalName: header.Filename,
		ShortCode:    shortCode,
		MimeType:     header.Header.Get("Content-Type"),
		Size:         header.Size,
		CreatedAt:    time.Now(),
		ExpiresAt:    expiryDate,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.Files.InsertOne(ctx, fileRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record file in database"})
		return
	}

	// Return success with the shortcode
	downloadURL := fmt.Sprintf("%s/download/%s", config.BaseURL, shortCode)
	c.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"shortCode":   shortCode,
		"downloadUrl": downloadURL,
		"expiresAt":   expiryDate,
	})
}

// DownloadFile handles file download requests
func DownloadFile(c *gin.Context) {
	shortCode := c.Param("shortcode")
	log.Printf("Download requested for shortcode: %s", shortCode)

	// Find file in database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var fileRecord database.FileRecord
	err := database.Files.FindOne(ctx, bson.M{"shortCode": shortCode}).Decode(&fileRecord)
	if err != nil {
		log.Printf("File not found in database: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if file has expired
	if time.Now().After(fileRecord.ExpiresAt) {
		log.Printf("File %s has expired", shortCode)
		c.JSON(http.StatusGone, gin.H{"error": "File has expired and is no longer available"})
		return
	}

	// Get file from SFTP
	remotePath := filepath.Join(config.SFTPPath, fileRecord.FileName)
	log.Printf("Fetching file from SFTP path: %s", remotePath)

	reader, err := utils.DownloadFromSFTP(remotePath)
	if err != nil {
		log.Printf("Failed to retrieve file from SFTP: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer reader.Close()

	// Set headers for file download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileRecord.OriginalName))
	c.Header("Content-Type", fileRecord.MimeType)

	// Stream file to client
	c.DataFromReader(
		http.StatusOK,
		fileRecord.Size,
		fileRecord.MimeType,
		reader,
		nil,
	)
}

// CleanupFiles manually triggers file cleanup
func CleanupFiles(c *gin.Context) {
	count := AutoCleanupFiles()
	c.JSON(http.StatusOK, gin.H{
		"message": "Cleanup completed",
		"deleted": count,
	})
}

// AutoCleanupFiles removes expired files (called by cron)
func AutoCleanupFiles() int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find expired files
	filter := bson.M{"expiresAt": bson.M{"$lt": time.Now()}}
	cursor, err := database.Files.Find(ctx, filter)
	if err != nil {
		log.Printf("Error finding expired files: %v", err)
		return 0
	}
	defer cursor.Close(ctx)

	deletedCount := 0

	// Delete each file from SFTP and database
	for cursor.Next(ctx) {
		var file database.FileRecord
		if err := cursor.Decode(&file); err != nil {
			continue
		}

		remotePath := filepath.Join(config.SFTPPath, file.FileName)
		err := utils.DeleteFromSFTP(remotePath)

		// Delete from database even if SFTP delete fails
		_, dbErr := database.Files.DeleteOne(ctx, bson.M{"_id": file.ID})

		if err == nil && dbErr == nil {
			deletedCount++
		}
	}

	return deletedCount
}
