
# syntax=docker/dockerfile:1
# ref: https://docs.docker.com/language/golang/build-images/

FROM golang:1.21
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY proxy ./proxy
COPY server ./server
COPY cmd ./cmd

RUN mkdir bin
RUN mkdir fs
RUN echo "FS=\"app/fs\"" > .env

RUN CGO_ENABLED=0 GOOS=linux go build -o proxy_bin /app/cmd/proxy/main.go
RUN mv proxy_bin bin/proxy
EXPOSE 6060

ENTRYPOINT ["/app/bin/proxy"]
