FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Use a minimal alpine image for running
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/app .

# Add necessary certificates for HTTPS
RUN apk --no-cache add ca-certificates

# The port will be provided by Render's environment
ENV PORT=8080

# Run the application
CMD ["/app/app"]
