# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.22.4@sha256:c8736b8dbf2b12c98bb0eeed91eef58ecef52b8c2bd49b8044531e8d8d8d58e8 AS build
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