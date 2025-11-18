FROM golang:1.21-alpine AS builder

# Install FFmpeg
RUN apk add --no-cache ffmpeg

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o /video-processor cmd/server/main.go

# Final stage
FROM alpine:latest

# Install FFmpeg in final image
RUN apk add --no-cache ffmpeg ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /video-processor .

# Create storage directories
RUN mkdir -p storage/uploads storage/processed storage/thumbnails

# Expose port
EXPOSE 8080

# Run application
CMD ["./video-processor"]
