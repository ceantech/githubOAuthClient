# syntax=docker/dockerfile:1

# Start from a base Golang image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /cmd

# Copy the Go module dependency files
# COPY cmd/main.go /cmd
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .
WORKDIR cmd

# Build the Go binary
RUN go build -o ghClient .

# Set the entrypoint for the container
ENTRYPOINT ["./ghClient"]

EXPOSE 9080
