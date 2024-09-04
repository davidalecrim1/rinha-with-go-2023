FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o server -ldflags="-w -s" ./cmd/main.go 

FROM scratch
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/database/schema.sql database/schema.sql
CMD ["./server"]