# syntax=docker/dockerfile:1
# ref: https://docs.docker.com/language/golang/build-images/

FROM golang:1.21
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
COPY server ./server
RUN CGO_ENABLED=0 GOOS=linux go build -o /http-server
EXPOSE 8080

ENTRYPOINT ["/http-server"]
CMD ["host-arg", "port-arg"]
