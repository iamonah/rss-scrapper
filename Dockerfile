# === Build Stage ===
FROM golang:1.23.2 AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# Copy the application source code
COPY . .

# Build the binary
RUN go build -ldflags="-s -w" -o /app/bin/main ./cmd/api/...

# === Runtime Stage ===
FROM gcr.io/distroless/static:nonroot

# Set working directory in the container
WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/bin/main /app/main

# Expose application port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
