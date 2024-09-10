# Start from the official golang image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o watchtower-line-notifier .

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/watchtower-line-notifier .

# Copy the example config file
COPY config.yaml.example ./config.yaml.example

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./watchtower-line-notifier"]