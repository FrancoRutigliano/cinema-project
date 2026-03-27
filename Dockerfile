# Build Stage
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o app ./cmd/

# Final Stage
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["/app/app"]
