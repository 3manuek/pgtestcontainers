package main

// XXX Testcontainers requires a self-hosted runner for
// execute docker-in-docker

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/3manuek/pgtestcontainers/internal/models"
	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
	postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	wait "github.com/testcontainers/testcontainers-go/wait"
)

func TestTimescale(t *testing.T) {
	ctx := context.Background()
	dbName := "iot"
	dbUser := "postgres"
	dbPassword := "password"

	usageData := filepath.Join("test/containerdata", "devices.csv")
	r, err := os.Open(usageData)
	if err != nil {
		t.Fatal(err)
	}
	postgresContainer, err := postgres.RunContainer(ctx,
		tc.WithImage("timescale/timescaledb:latest-pg16"),
		// We execute the generator with docker-compose, so we have a deterministic test
		// postgres.WithInitScripts(filepath.Join("test/containerdata", "003_generator.sql")),
		postgres.WithInitScripts(filepath.Join("test/timescale", "004_init.sql")),
		postgres.WithInitScripts(filepath.Join("test/timescale", "005_load.sql")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		tc.CustomizeRequest(tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Files: []tc.ContainerFile{
					{
						Reader:            r,
						HostFilePath:      usageData,
						ContainerFilePath: "/tmp/devices.csv",
						FileMode:          0o666,
					},
				},
			}}),
		tc.WithEnv(map[string]string{
			"TS_TUNE_MEMORY":   "1GB",
			"TS_TUNE_WAL":      "1GB",
			"TS_TUNE_NUM_CPUS": "2"}),
		tc.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second)), // we add a large startup due that we are loading data
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// This sleep handles the wait on processing files, we should do an active check here
	// instead of waiting
	time.Sleep(20 * time.Second)

	log.Println(postgresContainer.State(ctx))

	// Database pointer creation
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

	// add timeout to testCtx context for tests
	testCtx, testCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer testCtxCancel()

	// Test Case: Create NewDevices
	d, err := models.NewDevices(testCtx, db)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}

	// Add assertion that d.Bounds is not nil
	if d.Bounds == nil {
		t.Fatalf("Bounds are empty: %v", d.Bounds)
	}
	log.Println("Bounds:", d.Bounds)

	// Test Case: Aggregate without CAMV
	agg, _ := d.AggregateMaxCPUMinuteWOMV(testCtx)
	if len(agg) == 0 {
		t.Fatalf("AggregateMaxCPUMinuteWOMV is empty: %v", agg)
	}
	log.Println(len(agg))

	// Test Case: Execute plain command inside the container
	if _, out, err := postgresContainer.Exec(ctx, []string{"psql", "-U", dbUser, "-w", dbName, "-c", `SELECT count(*) from devices;`}); err != nil {
		log.Println(err)
		t.Fatal("couldn't count devices")
	} else {
		// read io.Reader out
		io.Copy(os.Stdout, out)
	}

	// Clean up the container
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
}
