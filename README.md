# Laboration 1: HTTP server

[Lab instructions](https://chalmers.instructure.com/courses/26458/pages/lab-1-http-server)

## Prereqs

1. Install [golang](https://go.dev/doc/install). Use go 1.21.x.
1. Install [docker](https://www.docker.com/get-started).

## Features

- Simple HTTP server ([server]('/server')) serving static files via GET and POST
- Simple proxy ([proxy]('/proxy')) implementation supporting GET requests

## Running

### Server

You can run the server locally or using docker.

#### Docker

```bash
docker build --rm --tag http_server -f server.Dockerfile .
docker run -p 8080:<port> http_server <ip> <port>
```

#### Locally

```bash
cd cmd/server
go build -o http_server main.go
./http_server <ip> <port>
```

### Proxy

#### Docker

```bash
docker build --rm --tag proxy -f proxy.Dockerfile .
docker run -p 6060:<port> proxy <port>
```

#### Locally

```bash
cd cmd/proxy
go build -o proxy main.go
./proxy <port>
```
