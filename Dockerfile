# First stage: build the Go binary
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o main ./cmd/app/main.go && chmod +x ./main


# Second stage: create the runtime image
FROM alpine:3.20.3

WORKDIR /app

# Install tzdata package for timezone support
RUN apk add --no-cache tzdata

# Set timezone to Asia/Jakarta
ENV TZ=Asia/Jakarta

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# # Copy the swagger.json file
# COPY --from=builder /app/docs/swagger.json ./docs/swagger.json

# Command to run the executable
CMD ["./main"]
