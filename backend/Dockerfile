# Stage 1: Build the Go binary
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run the binary
FROM alpine:latest
WORKDIR /app/

# Install dependencies required for wait-for-it
RUN apk add --no-cache bash

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Download the wait-for-it script
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it.sh

# Make the wait-for-it script executable
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Run the wait-for-it script before starting the main application
CMD ["wait-for-it.sh", "mysql:3306", "-t", "60", "--", "./main"]