# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM scratch

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .
# Copy CA certificates from builder (Alpine) stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set SSL_CERT_FILE environment variable (optional, but helps Go apps find the CA bundle)
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

# Run the application
CMD ["./main"]