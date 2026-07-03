# CLI Todo App

A simple todo application with:
- REST API server
- Postgres persistence
- CLI client

## Run server
go run ./cmd/server

## Run client
go run ./cmd/client list
go run ./cmd/client create --title "example"
