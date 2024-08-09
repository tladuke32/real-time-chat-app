# Stage 1: Build the Go binary
FROM golang:1.22.4-alpine AS builder
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
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Run the binary
CMD ["./main"]
