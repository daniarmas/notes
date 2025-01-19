# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.23.4@sha256:a6927f462c29ef917d9de1621c8e9ca5286948da4ea770f51f835a56c70cabc3 AS build
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notes -gcflags "all=-N -l" main.go

##
## Deploy
##
FROM gcr.io/distroless/static-debian11@sha256:63ebe035fbdd056ed682e6a87b286d07d3f05f12cb46f26b2b44fc10fc4a59ed

WORKDIR /app

COPY --from=build /app/notes /app/notes

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/notes", "run"]