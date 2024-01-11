package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	_ "github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/docs"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeletePhoto godoc
// @Summary Delete a photo
// @Description Deletes a photo for a given user and image ID
// @Tags photos
// @Produce  json
// @Param user_id path string true "User ID"
// @Param image_id path string true "Image ID"
// @Success 200 {object} string
// @Failure 404 {object} string
// @Router /deletephoto/{user_id}/{image_id} [delete]
func DeletePhoto(c *gin.Context) {
	startTimer := time.Now()
	// Extract the image_id and user_id from the request parameters
	imageID := c.Param("image_id")
	userID := c.Param("user_id")

	if _, err := uuid.Parse(imageID); err != nil {
		log.Printf("Invalid image ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID format"})
		return
	}

	if _, err := uuid.Parse(userID); err != nil {
		log.Printf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Prepare the delete statement
	db := database.NewDBConnection()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// First, delete related photo metadata
	_, err = tx.Exec("DELETE FROM photo_metadata WHERE image_id = $1", imageID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo metadata"})
		return
	}

	// Next, delete the photo
	_, err = tx.Exec("DELETE FROM images WHERE image_id = $1 AND user_id = $2", imageID, userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the photo"})
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	duration := time.Since(startTimer)
	prometheus.HTTPRequestDuration.WithLabelValues("/deletephoto").Observe(duration.Seconds())

	// If everything went well, send a success response
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
