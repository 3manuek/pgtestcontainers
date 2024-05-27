#!/bin/bash

docker compose up -d 

sleep 30

go run main.go

docker compose down -v

go test -v compose_test.go

go test -v generic_test.go

go test -v ts_test.go
