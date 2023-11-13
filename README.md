# Laboration 1: HTTP server

[Lab instructions](https://chalmers.instructure.com/courses/26458/pages/lab-1-http-server)

## Prereqs

1. Install [golang](https://go.dev/doc/install). Use go 1.21.x.
1. Install [docker](https://www.docker.com/get-started).

## Features

- Simple HTTP server ([server]('/server')) serving static files via GET and POST
- Simple proxy ([proxy]('/proxy')) implementation supporting GET requests

## Cloud

- The HTTP server is hosted on AWS and can be accessed at `52.203.125.74:80`
- The proxy is hosted on AWS and can be accessed at `54.81.33.35:80`

## Running

### Setting the environment variable

Our solution uses godotenv to set a path to the directory which is used for storing and getting files, you need to specify a path to make it work.

```bash
FS="/set/your/abs/path/here"
```

### Tests

The project has been tested on MacOS and Fedora Linux.

```
go test ./... -v
```

### Server

You can run the server locally or using docker.

#### Docker

```bash
docker build --rm --tag http_server -f server.Dockerfile .
docker run -p 80:<port> http_server
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
docker run -p 80:<port> proxy 
```

#### Locally

```bash
cd cmd/proxy
go build -o proxy main.go
./proxy <port>
```
