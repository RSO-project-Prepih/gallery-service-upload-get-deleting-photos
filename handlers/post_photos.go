package handlers

import (
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadPhoto(c *gin.Context) {
	startTime := time.Now()

	var input struct {
		Username  string `form:"username" binding:"required"`
		UserID    string `form:"user_id" binding:"required"`
		ImageName string `form:"image_name" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileInterface, exists := c.Get("file")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File processing failed"})
		return
	}
	file := fileInterface.(*multipart.FileHeader)

	// Read the file's content
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file"})
		return
	}
	defer fileContent.Close()

	// Convert the file's content to bytes
	bytesContent := make([]byte, file.Size)
	_, err = fileContent.Read(bytesContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file content"})
		return
	}

	// Insert the file's content into the database
	db := database.NewDBConnection()
	defer db.Close()
	result, err := db.Exec(
		"INSERT INTO images (image_name, data, image_id, user_id) VALUES ($1, $2, $3, $4)",
		input.ImageName,
		bytesContent,
		uuid.New().String(),
		input.UserID,
	)
	if err != nil {
		log.Printf("Failed to save the photo to the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the photo to the database"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to retrieve affected rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows"})
		return
	}

	log.Printf("Photo uploaded successfully. Rows affected: %d", rowsAffected)

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded successfully"})
	duration := time.Since(startTime)
	prometheus.HTTPRequestDuration.WithLabelValues("/uploadphoto").Observe(duration.Seconds())
}
