FROM golang:1.23-alpine AS builder

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

# libc6-compat is often needed for CGO binaries on alpine
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /root/

# Create and set permissions for the uploads and data directories
RUN mkdir -p /root/uploads /root/data && \
    chmod 777 /root/uploads /root/data

COPY --from=builder /app/main .
COPY --from=builder /app/internal/database/migrations ./migrations

EXPOSE 8080

# DB_URL will point to the SQLite file in /root/data/
ENV DB_URL=/root/data/skillswap.db

CMD ["./main"]
