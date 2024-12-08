# Use an Alpine base image with Golang
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Install necessary dependencies (e.g., bash, git)
RUN apk add --no-cache bash git

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/main.go

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
