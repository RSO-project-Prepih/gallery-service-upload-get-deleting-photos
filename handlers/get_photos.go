package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PhotoResponse represents a photo object returned by the API
type PhotoResponse struct {
	ImageID   string `json:"image_id"`
	ImageName string `json:"image_name"`
	UserID    string `json:"user_id"`
	Data      string `json:"data"`
}

// DisplayPhotos retrieves all photos for a specific user by their ID
func DisplayPhotos(c *gin.Context) {
	userID := c.Param("user_id")

	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Query the database for all images belonging to the user
	rows, err := database.DB.Query("SELECT image_id, image_name, user_id, data FROM images WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query photos"})
		return
	}
	defer rows.Close()

	// Iterate over the rows and construct a slice of Photos\

	photos := []PhotoResponse{}
	for rows.Next() {
		var (
			imageID, imageName, uid string
			data                    []byte
		)
		if err := rows.Scan(&imageID, &imageName, &uid, &data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan photos"})
			return
		}

		// Encode the image data to base64
		encodedData := base64.StdEncoding.EncodeToString(data)

		photos = append(photos, PhotoResponse{
			ImageID:   imageID,
			ImageName: imageName,
			UserID:    uid,
			Data:      encodedData,
		})
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over photos"})
		return
	}

	c.JSON(http.StatusOK, photos)
}
