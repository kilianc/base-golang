# 1.21.5 golang image
FROM golang:1.21.5-alpine3.18

WORKDIR /app
COPY [".", "./"]

RUN go mod download
