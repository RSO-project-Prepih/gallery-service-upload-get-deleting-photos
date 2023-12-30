package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/RSO-project-Prepih/gallery-service-uplode-get-deliting-photos.git/database"
	"github.com/gin-gonic/gin"
)

// PhotoMetadata struct as before
type PhotoMetadata struct {
	Id           int             `json:"id"`
	ImageId      string          `json:"image_id"`
	CameraModel  sql.NullString  `json:"camera_model,omitempty"`
	Date         sql.NullString  `json:"date,omitempty"`
	Latitude     sql.NullFloat64 `json:"latitude,omitempty"`
	Longitude    sql.NullFloat64 `json:"longitude,omitempty"`
	Altitude     sql.NullFloat64 `json:"altitude,omitempty"`
	ExposureTime sql.NullString  `json:"exposure_time,omitempty"`
	LensModel    sql.NullString  `json:"lens_model,omitempty"`
	CreatedAt    sql.NullString  `json:"created_at,omitempty"`
	UpdatedAt    sql.NullString  `json:"updated_at,omitempty"`
}

func GetPhotoMetadataByUser(c *gin.Context) {
	userID := c.Param("user_id")

	db := database.NewDBConnection()
	defer db.Close()

	var metadata []PhotoMetadata

	query := `
        SELECT pm.id, pm.image_id, pm.camera_model, pm.date, pm.latitude, pm.longitude, pm.altitude, pm.exposure_time, pm.lens_model, pm.created_at, pm.updated_at 
        FROM photo_metadata pm
        INNER JOIN images i ON pm.image_id = i.image_id
        WHERE i.user_id = $1
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query photo metadata"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var m PhotoMetadata
		err := rows.Scan(&m.Id, &m.ImageId, &m.CameraModel, &m.Date, &m.Latitude, &m.Longitude, &m.Altitude, &m.ExposureTime, &m.LensModel, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		metadata = append(metadata, m)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate photo metadata"})
		return
	}

	c.JSON(http.StatusOK, metadata)
}
