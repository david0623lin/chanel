# Base image
FROM golang:1.20.5-alpine

# code review image
FROM golangci/golangci-lint:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Add Go binary path to the environment variables
ENV PATH=$HOME/go/bin:$PATH

# Initialize swag
RUN swag init

# run code review (timeout 10m)
RUN golangci-lint run --timeout 10m

# Build the Go app
RUN go build -mod=vendor -o main .

# Start command
CMD ["./main"]
