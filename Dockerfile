# Start from the official Golang image
FROM golang:1.23.2-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go app (output binary name is benchmark-api)
RUN go build -o benchmark-api

# Use a smaller image for running the app
FROM alpine:3.16

# Install ca-certificates for making HTTP requests
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder image
COPY --from=builder /app/benchmark-api /usr/local/bin/benchmark-api

# Expose port 8080 for the Gin API
EXPOSE 8080

# Run the benchmark app
CMD ["benchmark-api"]  # Ensure the binary name matches the one built
