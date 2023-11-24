
# Use the official Golang image for the application
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Install any dependencies if needed
RUN go mod download

# Build the Go application
RUN go build -o main ./cmd/api

# Expose the application port
EXPOSE 4000

# Define the command to run the application
CMD ["./main"]