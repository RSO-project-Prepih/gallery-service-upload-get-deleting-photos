package main

import (
	"log"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server...")

	r := gin.Default()

	r.POST("/upload", handlers.UploadPhoto)
	//r.GET("/photos/:user_id", handlers.FetchPhotos)
	//r.DELETE("/photo/:photo_id", handlers.DeletePhoto)

	r.Run(":8080")

}
