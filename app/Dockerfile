# Start from the official Go image
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod tidy

# Download dependencies
RUN go mod download

# Copy the entire Go app source code into the container
COPY . .

# Build statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix 'static' -o go-api .

# Stage 2 — Minimal Chainguard static image
FROM cgr.dev/chainguard/static:latest

WORKDIR /app

# Copy the built executable from the builder image
COPY --from=builder /app/go-api .

# Expose the port your app will run on (e.g., port 8080)
EXPOSE 8080

# Set the entrypoint to your Go app
ENTRYPOINT ["/app/go-api"]