package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/handlers"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting the application...")
	r := gin.Default()

	r.POST("/uploadphotos", middlewares.CheckImageContentType(), handlers.UploadPhoto)
	r.GET("/getphotos/:user_id", handlers.DisplayPhotos)
	r.DELETE("/deletephoto/:user_id/:image_id", handlers.DeletePhoto)

	srver := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	qoit := make(chan os.Signal, 1)
	signal.Notify(qoit, syscall.SIGINT, syscall.SIGTERM)
	<-qoit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srver.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
