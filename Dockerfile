# Build stage
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./

FROM alpine:latest

RUN apk add --no-cache \
    wget \
    unzip 

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080 4444

CMD ["./main"]