package handlers

import (
	"net/http"
	"time"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeletePhoto(c *gin.Context) {
	startTimer := time.Now()
	// Extract the image_id and user_id from the request parameters
	imageID := c.Param("image_id")
	userID := c.Param("user_id")

	if _, err := uuid.Parse(imageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID format"})
		return
	}

	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Prepare the delete statement
	db := database.NewDBConnection()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM images WHERE image_id = $1 AND user_id = $2")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare the delete statement"})
		return
	}
	defer stmt.Close()

	// Execute the delete statement
	result, err := stmt.Exec(imageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the photo"})
		return
	}

	// Check if the photo was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check the affected rows"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No photo found with the specified ID for the user"})
		return
	}
	duration := time.Since(startTimer)
	prometheus.HTTPRequestDuration.WithLabelValues("/deletephoto").Observe(duration.Seconds())

	// If everything went well, send a success response
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
