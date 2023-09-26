# Use the official Go image from the DockerHub
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files to the app directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o member-counts-go .

# Start fresh from a smaller image
FROM alpine:latest

# Install runtime dependencies and create app user
RUN apk --no-cache add ca-certificates tzdata && addgroup -S app && adduser -S -G app app

# Security related setting, switch to non-root user
USER app

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/member-counts-go .

# Command to run the executable
ENTRYPOINT ["./member-counts-go"]
