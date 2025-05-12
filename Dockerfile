# Build stage. "Install" Golang inside the container
FROM golang:1.24 AS builder

# directory inside the container
# copies code from src to app
# build the go app

WORKDIR /app          
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o server main.go

# Runtime image
FROM debian:bullseye-slim

# Install CA certificates for scrapper to work properly
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# copy compiled binary to image
WORKDIR /app
COPY --from=builder /app/server .

# Start!
EXPOSE 8080
ENTRYPOINT ["/app/server"]