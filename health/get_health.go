package health

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/joho/godotenv"
)

func GetDSN() string {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		// You may want to handle this error differently than a fatal log, depending on your application's requirements
		return ""
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func HealthCheckHandler(db *sql.DB) http.Handler {
	// Create a new health check handler
	health := healthcheck.NewHandler()

	// Add a readiness check to verify database connectivity
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 1*time.Second))

	// Add other checks here, such as liveness checks, if needed

	// Return the handler
	return health
}
