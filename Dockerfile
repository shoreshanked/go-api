# Start from the official Go image
FROM golang:1.21 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod tidy

# Download dependencies
RUN go mod download

# Copy the entire Go app source code into the container
COPY . .

# Build the Go application
RUN go build -o go-api .

# Final image (scratch or lightweight base)
FROM ubuntu:22.04

WORKDIR /app

# Install tzdata package to handle timezones
RUN apt-get update && apt-get install -y tzdata

# Copy the built executable from the builder image
COPY --from=builder /app/go-api .

# Expose the port your app will run on (e.g., port 8080)
EXPOSE 8080

# Run the Go application
#CMD ["./go-api"]
CMD ["sh", "-c", "while true; do sleep 1000; done"]