# Use an official Go runtime as a parent image
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY go-src/. .

# Download Go modules and build the Go application
RUN go mod download & go build -o main .

# Use a minimal Alpine Linux image for the final runtime
FROM alpine:3.23.3

# Set the working directory
WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /app/main .

# Expose any ports your application uses (if applicable)
# EXPOSE 8080

# Command to run the application
CMD ["./main"]
