# syntax=docker/dockerfile:1
# ref: https://docs.docker.com/language/golang/build-images/

FROM golang:1.21
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY server ./server
COPY cmd ./cmd

RUN mkdir fs
RUN echo "FS=\"app/fs\"" > .env
RUN chmod +x /app/server

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/http_server /app/cmd/server/main.go
EXPOSE 8080

ENTRYPOINT ["/app/http_server"]
