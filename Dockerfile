# Use the official Go image from the DockerHub
FROM golang:1.20

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

# Security related setting, see https://stackoverflow.com/questions/52215283/what-is-the-use-of-user-nobody-in-dockerfile
RUN apk --no-cache add ca-certificates && addgroup -S app && adduser -S -G app app
USER app

# Copy the Pre-built binary file from the previous stage
COPY --from=0 /app/member-counts-go .

# Command to run the executable
ENTRYPOINT ["./member-counts-go"]
