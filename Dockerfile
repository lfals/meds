# Use the official Go image as the base image
FROM golang:1.23.4-alpine

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files (if they exist)
COPY go.mod ./
COPY go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application with CGO enabled
RUN CGO_ENABLED=1 go build -o main .

# Expose port 3000
EXPOSE 3000

# Command to run the executable
CMD ["./main"]
