FROM golang:1.22.2-alpine AS builder
WORKDIR /app

COPY go.sum go.mod ./
RUN go mod download
COPY .  ./
RUN go build -ldflags '-s -w' -o main .

# Use a lightweight image for final distribution
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main ./

CMD ["./main"]
