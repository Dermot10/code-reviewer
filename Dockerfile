# Use official Go image
FROM golang:1.21-alpine

# Install git (needed for module fetching)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod + go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your code
COPY . .

# Build your Go binary
RUN go build -o server ./cmd/setup.go

# Expose port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./server"]
