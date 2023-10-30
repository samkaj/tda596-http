#!/bin/bash

docker build --rm --tag http-server .

# TODO: when server is listening, uncomment below
# docker run -p 8080:8080 http-server

docker run http-server # TODO: and remove this line
