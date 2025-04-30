# Use the same Go version and Alpine version as KrakenD to ensure compatibility
FROM golang:1.23.7-alpine3.21 AS builder

# Install necessary tools and dependencies
RUN apk add --no-cache musl-dev gcc git

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Ensure required dependencies are fetched and aligned with KrakenD
RUN go clean -modcache
RUN go mod tidy
RUN go mod download

# Copy the source code into the container
COPY validator.go ./

# Set environment variables for plugin compilation
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Build the plugin(Convert .go file to .so file which would be utilised by krakend service)
RUN go build -buildmode=plugin -o validator.so validator.go

# Second stage: Use the KrakenD runtime image
FROM devopsfaith/krakend:2.9.3

# Set the working directory for KrakenD
WORKDIR /app

# Copy the compiled plugin and dependency files from the builder stage
COPY --from=builder /app/validator.so /app/
COPY --from=builder /app/go.mod /app/
COPY --from=builder /app/go.sum /app/

# Expose the default KrakenD port
EXPOSE 3000

# Default command for validating the plugin
CMD ["check-plugin", "/app/validator.so"]