package main

import (
	"context"
	"flag"
	"fmt"
	"testing"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// go test -v main_test.go -args -pgaURL=...
var imageName = flag.String("imageName", "postgres:16-bookworm", "URL of the image")

// var registry = flag.String("registry", "", "Registry domain")

// TODO implement https://www.lambdatest.com/automation-testing-advisor/golang/methods/testcontainers-go_go.wait.ForSQL
// https://golang.testcontainers.org/quickstart/
func TestGenericContainer(t *testing.T) {
	// var URL = *registry + "/" + *imageName
	ctx := context.Background()
	req := tc.ContainerRequest{
		Image:        *imageName,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD":         "postgres",
			"POSTGRES_HOST_AUTH_METHOD": "trust"},
		WaitingFor: wait.ForLog("Ready to accept connections"),
	}
	postgresC, _ := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	defer func() {
		if err := postgresC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	state, err := postgresC.State(ctx)
	if err != nil {
		t.Fatalf("failed to get container state: %s", err) // nolint:gocritic
	}

	fmt.Println(state.Running)
	if _, _, err := postgresC.Exec(ctx, []string{"psql", "-U", "postgres", "-P", "postgres", "-c", "SELECT version();"}); err != nil {
		t.Fatal(err)
	}

}
