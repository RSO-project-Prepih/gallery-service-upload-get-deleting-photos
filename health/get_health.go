package health

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/heptiolabs/healthcheck"
	_ "github.com/lib/pq"
)

func GetDSN() string {

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=verify-full",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func HealthCheckHandler(db *sql.DB) (http.HandlerFunc, http.HandlerFunc) {

	// Create a health instance.
	health := healthcheck.NewHandler()

	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	// Add a readiness check for a database.
	health.AddReadinessCheck(
		"database",
		healthcheck.DatabasePingCheck(db, 1*time.Second),
	)

	return health.LiveEndpoint, health.ReadyEndpoint
}
