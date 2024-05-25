# pgcontainers

A repository with examples for testing Postgres containers with 
Testcontainers.

```bash
docker compose up -d
go run main.go
```

## Postgres with Generic Container Requests

```bash
go test generic_test.go
```

## Postgres (Timescale version) with Postgres Module

````bash
go test ts_test.go
```

