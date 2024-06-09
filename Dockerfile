# Start from the latest golang base image
FROM golang:1.21.5 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .


# List the contents of the /app/ directory
RUN ls -la /app/

# Build the Go app for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Start a new stage from scratch
FROM ubuntu:latest  

WORKDIR /root/

# Install curl
RUN apt-get update && apt-get install -y curl

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the script from the previous stage
COPY --from=builder /app/scripts/startup.sh .

# List the contents of the /root/ directory
RUN ls -la /root/

# Make your script executable
RUN chmod +x ./startup.sh

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["sh", "-c", "./startup.sh && ./main"]