
FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Create and set permissions for the uploads directory
RUN mkdir -p /root/uploads && \
    chmod 777 /root/uploads

COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads
COPY --from=builder /app/database/migrations ./migrations

EXPOSE 8080

CMD ["./main"]
