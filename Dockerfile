FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY /cmd ./cmd 
COPY /config ./config 
COPY /database ./database
COPY /internal ./internal

RUN CGO_ENABLED=0 go build -o server ./cmd/main.go 

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/database/schema.sql database/schema.sql
CMD ["./server"]