# multitenant-api-go

Rest API for a multitenant system using golang. It includes a Rest API using **GIN** and **Postgres** for the database.

## Components

- *Database*: It uses **Postgres** for the database
- *Documentation*: Swagger API documentation

## Requirements

- [golang 1.22](https://go.dev/doc/install)

### optional

- Run `make install` for optional dependencies (air and swag)

## Development

- Execute `go run ./cmd/server/main.go` to start the server.
- Makefile included to simplify commands:
  - `make dev` start de server on live reload mode using *air*
  - `make run` to start the server
  - `make get` to get dependencies

- Alternatively, you can use docker and docker-compose to run the database and the server
  - `make dbuild` to build the *api* image
  - `make dup` to start the *api* container
  - `make dlogs` to visualize the *api* container logs

## Packages

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM (Golang ORM)](https://gorm.io/)
- [air](https://github.com/cosmtrek/air)
- [swag](https://github.com/swaggo/swag)
