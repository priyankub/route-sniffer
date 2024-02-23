# Use the official Go image as the base image
FROM golang:alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o /app/route-sniffer ./cmd

# Use a lightweight Alpine base image to create the final image
FROM alpine

# Copy the built executable from the previous stage
COPY --from=builder /app/route-sniffer /app/route-sniffer

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the Go web application
CMD ["/app/route-sniffer"]
