# Build stage
FROM --platform=linux/amd64 golang:1.20 AS builder

# Set up the working directory
WORKDIR /app/src

# Fetch dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/bin/main .

# Runtime stage
FROM --platform=linux/amd64 alpine:3.14

# Copy the binary from the builder stage
COPY --from=builder /app/bin/main /app/bin/main

RUN chmod +x /app/bin/main

# Set the binary as the entrypoint
ENTRYPOINT ["/app/bin/main"]