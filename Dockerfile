FROM golang:1.22.2-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

#Download dependency
RUN go mod tidy

COPY . .

# Build the binary and set executable permission in one step
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o main ./cmd/app/main.go && chmod +x ./main


# Second stage: create the runtime image
FROM alpine:3.20.3

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the swagger.json file
COPY --from=builder /app/docs/swagger.json ./docs/swagger.json

#Try adding environment Variable
ENV TZ=Asia/Jakarta

# Command to run the executable
CMD ["./main"]                    