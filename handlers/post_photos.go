package handlers

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	get_photo_info "github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/github.com/RSO-project-Prepih/get-photo-info"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/grpcclient"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadPhoto(c *gin.Context) {
	startTime := time.Now()
	log.Printf("UploadPhoto handler started at %v", startTime)

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

	imageID := uuid.New().String()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photo info"})
		return
	}

	// Insert the file's content into the database
	db := database.NewDBConnection()
	defer db.Close()
	result, err := db.Exec(
		"INSERT INTO images (image_name, data, image_id, user_id) VALUES ($1, $2, $3, $4)",
		input.ImageName,
		bytesContent,
		imageID,
		input.UserID,
	)
	if err != nil {
		log.Printf("Failed to save the photo to the database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the photo to the database"})
		return
	}
	log.Printf("Photo saved to the database successfully")
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to retrieve affected rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows"})
		return
	}
	log.Printf("Rows affected: %d", rowsAffected)

	grpcAddress := os.Getenv("GRPC_ADDRESS")
	if grpcAddress == "" {
		log.Fatal("GRPC_ADDRESS environment variable not set")
	}
	photoServiceClient, err := grpcclient.NewPhotoServiceClient(grpcAddress)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo service client"})
		return
	}

	log.Printf("Photo service client created successfully")
	log.Printf("photoServiceClient: %v", photoServiceClient)

	response, err := photoServiceClient.GetPhotoInfo(context.Background(), &get_photo_info.PhotoRequest{
		Photo:   bytesContent,
		ImageId: imageID,
	})

	log.Printf("Response: %v", response)
	if response == nil {
		log.Printf("Failed to get response from the photo service: %v", err)
	}

	if err != nil || !response.Allowed {
		// Delete the image if the response is not allowed
		_, delErr := db.Exec("DELETE FROM images WHERE image_id = $1", imageID)
		if delErr != nil {
			log.Printf("Failed to delete the photo: %v", delErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the photo"})
			return
		}
		log.Printf("Photo deleted successfully")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Photo not allowed"})
		return
	}

	log.Printf("Photo uploaded successfully. Rows affected: %d", rowsAffected)

	c.JSON(http.StatusOK, gin.H{"message": "Photo uploaded successfully"})
	duration := time.Since(startTime)

	prometheus.HTTPRequestDuration.WithLabelValues("/uploadphoto").Observe(duration.Seconds())

}
