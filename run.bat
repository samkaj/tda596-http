@echo off

REM Build the Docker image with the tag "http-server"
docker build --rm --tag http-server .

REM TODO: When server is listening, uncomment below
REM docker run -p 8080:8080 http-server

REM Run the Docker container
docker run http-server REM TODO: and remove this line
