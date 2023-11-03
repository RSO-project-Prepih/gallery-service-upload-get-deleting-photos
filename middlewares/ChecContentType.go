package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckImageContentType() gin.HandlerFunc {
	return func(c *gin.Context) {

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			c.Abort()
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open the file"})
			c.Abort()
			return
		}
		defer openedFile.Close()

		buffer := make([]byte, 512)
		_, err = openedFile.Read(buffer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file content"})
			c.Abort()
			return
		}

		contentType := http.DetectContentType(buffer)
		if !strings.HasPrefix(contentType, "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The uploaded file is not an image"})
			c.Abort()
			return
		}

		// If the content type is correct, reset the read pointer and proceed
		_, err = openedFile.Seek(0, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file content"})
			c.Abort()
			return
		}

		c.Set("file", file) // Pass the file to the next handler
		c.Next()
	}
}
