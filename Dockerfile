FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-shortener ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/url-shortener .

EXPOSE 8080

CMD ["./url-shortener"]
