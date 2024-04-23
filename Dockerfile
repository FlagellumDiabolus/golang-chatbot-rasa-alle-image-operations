# Use an official Golang runtime as a parent image
FROM golang:1.17-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o golang-chatbot-alle-image-operations ./cmd/chatbot-api

# Start a new stage from scratch
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the current working directory
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/golang-chatbot-alle-image-operations .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./golang-chatbot-alle-image-operations"]
