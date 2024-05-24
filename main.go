package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"log"
	"time"

	"github.com/3manuek/pgtestcontainers/internal/models"
)

func main() {

	// Database pointer creation
	connStr := "postgres://postgres:postgres@localhost:15432/iot?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// add timeout to testCtx context for tests
	testCtx, testCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer testCtxCancel()

	// Test Case
	d, err := models.NewDevices(testCtx, db)
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	agg, _ := d.AggregateMaxCPUMinuteWOMV(testCtx)
	log.Println(len(agg))

	log.Printf("Bounds: %s, %s", (d.Bounds).StartTime.Format(time.RFC3339), (d.Bounds).EndTime.Format(time.RFC3339))
}
