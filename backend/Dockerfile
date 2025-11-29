FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server .


FROM alpine:latest
WORKDIR /app
# Copy the pre-built binary
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
# Expose port (from .env)
EXPOSE ${APP_PORT}
# Run the binary
CMD ["/app/server"]