package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/docs"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/handlers"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/health"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/middlewares"
	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gallery Service API
// @version 1.0
// @description This is a service for uploading, getting and deleting photos
// @BasePath /v1
func main() {
	log.Println("Starting the application...")
	r := gin.Default()

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour / time.Second),
	})

	r.Use(func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	})

	// handelers
	r.POST("/uploadphotos", middlewares.CheckImageContentType(), handlers.UploadPhoto)
	r.GET("/getphotos/:user_id", handlers.DisplayPhotos)
	r.DELETE("/deletephoto/:user_id/:image_id", handlers.DeletePhoto)
	r.GET("/photometadata/:user_id", handlers.GetPhotoMetadataByUser)

	// live and ready
	liveHandler, readyHandler := health.HealthCheckHandler()
	r.GET("/live", gin.WrapH(liveHandler))
	r.GET("/ready", gin.WrapH(readyHandler))

	// metrics
	r.GET("/metrics", gin.WrapH(prometheus.GetMetrics()))

	// swagger
	r.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// server
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
