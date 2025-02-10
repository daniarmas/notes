## About The Project

This is a backend application for a notes app, developed using Go. The API supports both REST and GraphQL communication. While the application is designed to be simple in terms of business logic, it serves as a robust example of a Go project structure. It is intended to be a solid starting point for more complex applications.

## Built With

* Programming language - [Go](https://go.dev)
* Http server - [net/http](https://pkg.go.dev/net/http)
* GraphQL library - [gqlgen](https://gqlgen.com)
* Database - [PostgreSQL](https://www.postgresql.org)
* Cache - [Redis](https://redis.io/)

## How to test

1. Install [Docker](https://docs.docker.com/desktop/)
2. Install Bruno from [here](https://www.usebruno.com/downloads).
3. Clone the repo
   ```sh
   git clone https://github.com/daniarmas/notes.git
   ```
4. Open the Bruno collections at **./api/bruno**.
5. Configure an object storage service compatible with the Amazon S3 API. [DigitalOcean Spaces](https://docs.digitalocean.com/products/spaces/) was used in the development. Ensure you update the ***access key***, ***secret key*** and ***bucket name*** in the docker compose file enviroment variables.
6. Create the `.env` file with the env vars in example.env
   ```sh
   cp example.env .env
   ```
7. Deploy docker compose file
   ```sh
   docker compose -f deploy/docker-compose.yaml up -d
   ```

## Setup for development

1. Install [Go 1.22.4](https://go.dev/doc/install)
2. Install [Docker](https://docs.docker.com/desktop/)
3. Install [direnv](https://direnv.net) to export the environment variables. (*only for development*)
4. Clone the repo
   ```sh
   git clone https://github.com/daniarmas/notes.git
   ```
5. Install Go dependencies
   ```sh
   go mod download
   ```
6. Deploy docker compose file
   ```sh
   docker compose -f deploy/docker-compose-dev.yaml up -d
   ```
7. Run direnv command to approve his content
   ```sh
   direnv allow
   ```
8. Create the `.envrc` file with the env vars in example.envrc
9. ```sh
   cp example.envrc .envrc
   ```
10. Create the database tables
   ```sh
   go run main.go create database
   ```
11. Seed the database. This seed the database for test purpose
   ```sh
   go run main.go create seed
   ```
12. Configure an object storage service compatible with the Amazon S3 API. [DigitalOcean Spaces](https://docs.digitalocean.com/products/spaces/) was used in the development. Ensure you update the ***access key***, ***secret key*** and ***bucket name*** in the `.envrc` file.
13.  Run the app
   ```sh
   go run main.go run
   ```

## Documentation

After completing step 3 in the [Setup for development](#setup-for-development) section, go to [http://localhost:8081](http://localhost:8081) to open Swagger UI.

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

## Contact

[Contact information](https://github.com/daniarmas)

## Acknowledgments

* [Sqlc](https://docs.sqlc.dev/en/latest/#)
* [PostgreSQL Driver](https://github.com/jackc/pgx)
* [Library for CLI apps](https://github.com/spf13/cobra)
* [Redis client](https://github.com/redis/go-redis/)