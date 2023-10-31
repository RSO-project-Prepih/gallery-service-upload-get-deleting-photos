package handlers

import (
	"fmt"
	"net/http"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadPhoto(c *gin.Context) {
	var input struct {
		Username  string `form:"username" binding:"required"`
		UserID    string `form:"user_id" binding:"required"`
		ImageName string `form:"image_name" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

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
	_, err = database.DB.Exec(
		"INSERT INTO images (image_name, data, image_id, user_id) VALUES ($1, $2, $3, $4)",
		input.ImageName,
		bytesContent,
		uuid.New().String(),
		input.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the photo to the database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded successfully"})

}

func FetchPhotos() {
	fmt.Println("FetchPhotos")
}

func DeletePhoto() {
	fmt.Println("DeletePhoto")
}
