# Start with the official Golang image to build the Go application
FROM golang:1.22.3 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o Envio

# Use a smaller base image for the final container
#FROM debian:buster-slim

# Stage 2: Create a small image using Alpine Linux
FROM alpine:3.14

# Install any necessary packages, for example, CA certificates if needed
# RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/Envio /usr/local/bin/Envio

# Command to run the executable
CMD ["/usr/local/bin/Envio"]
