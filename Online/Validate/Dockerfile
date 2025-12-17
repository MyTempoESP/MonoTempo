# Start with the official Golang image to build the Go application
FROM golang:1.22.3 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies files
COPY go.mod go.sum ./

COPY vendor/ .vendor/

# Copy the source code into the container
COPY . .

ENV GOFLAGS="-mod=vendor"

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o Reenvio

# Use a smaller base image for the final container
#FROM debian:buster-slim

# Stage 2: Create a small image using Alpine Linux
FROM alpine:3.14

# Install any necessary packages, for example, CA certificates if needed
# RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/Reenvio /usr/local/bin/Reenvio

# Command to run the executable
CMD ["/usr/local/bin/Reenvio"]
