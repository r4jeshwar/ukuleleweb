FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y curl && \
    curl -O https://dl.google.com/go/go1.20.1.linux-amd64.tar.gz && \
    tar -xvf go1.20.1.linux-amd64.tar.gz && \
    mv go /usr/local

ENV PATH="/usr/local/go/bin:${PATH}"

RUN bash

WORKDIR ukeleleweb

COPY . .

RUN sed "s/localhost/0.0.0.0/g" cmd/ukuleleweb/main.go -i

RUN cd cmd/ukuleleweb/ && \
    go mod tidy && \
    go build && \
    mkdir ukuleleweb-data
    
EXPOSE 8080

ENTRYPOINT ["./cmd/ukuleleweb/ukuleleweb", "-store_dir=ukuleleweb-data"]

# # First stage: build the Go application
# FROM golang:1.20.1-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN CGO_ENABLED=0 go build -o app

# # Second stage: create the final container
# FROM alpine:latest

# WORKDIR /app

# COPY --from=builder /app/app ./

# CMD ["./app"]
