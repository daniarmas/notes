<!-- ABOUT THE PROJECT -->
## About The Project

This is a backend application for a notes app, developed using Go. The API supports both REST and GraphQL communication. While the application is designed to be simple in terms of business logic, it serves as a robust example of a Go project structure. It is intended to be a solid starting point for more complex applications.

## Built With

* Programming language - [Go](https://go.dev)
* Http server - [net/http](https://pkg.go.dev/net/http)
* GraphQL library - [gqlgen](https://gqlgen.com)
* Database - [PostgreSQL](https://www.postgresql.org)
* Cache - [Redis](https://redis.io/)

## How to test

1. Install Postman from [here](https://www.postman.com/downloads/).
2. Download the api postman collection [here](https://github.com/daniarmas/notes/blob/main/assets/notes-api.postman_collection.json)

## Documentation

After completing step 3 in the [Setup for development](#setup-for-development) section, go to [http://localhost:8081](http://localhost:8081) to open Swagger UI.

<!-- PREREQUISITES -->
## Prerequisites

1. Install [Go 1.22.4](https://go.dev/doc/install)
2. Install [Docker](https://docs.docker.com/desktop/)
3. Install [direnv](https://direnv.net) to export the environment variables. (*only for development*)

<!-- INSTALLATION -->
## Setup for development

1. Clone the repo
   ```sh
   git clone https://github.com/daniarmas/notes.git
   ```
2. Install Go dependencies
   ```sh
   go mod download
   ```
3. Deploy docker compose file
   ```sh
   docker compose -f deploy/docker-compose-dev.yaml up -d
   ```
4. Run direnv command to approve his content
   ```sh
   direnv allow
   ```
5. Create the `.envrc` file with the env vars in example.envrc
6. ```sh
   cp example.envrc .envrc
   ```
7. Create the database tables
   ```sh
   go run main.go create database
   ```
8. Seed the database. This seed the database for test purpose
   ```sh
   go run main.go create seed
   ```
9. Configure an object storage service compatible with the Amazon S3 API. [DigitalOcean Spaces](https://docs.digitalocean.com/products/spaces/) was used in the development. Ensure you update the ***access key***, ***secret key*** and ***bucket name*** in the `.envrc` file.
10.  Run the app
   ```sh
   go run main.go run
   ```

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<!-- CONTACT -->
## Contact

[Contact information](https://github.com/daniarmas)

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Sqlc](https://docs.sqlc.dev/en/latest/#)
* [PostgreSQL Driver](https://github.com/jackc/pgx)
* [Library for CLI apps](https://github.com/spf13/cobra)
* [Redis client](https://github.com/redis/go-redis/)