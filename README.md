# pgcontainers

A repository with examples for testing Postgres containers with 
Testcontainers in a variety of ways.

```bash
docker compose up -d
go run main.go
```

## Postgres with Generic Container Requests

```bash
go test generic_test.go
```

## Postgres (Timescale version) with Postgres Module

```bash
go test ts_test.go
```

## Todo

- [ ] Implement a wait.ForSQL method for checking specific rules against the model.
- [ ] Implement a wait.ForDBRule for testing a expected value in the database, with a timeout'ed context.