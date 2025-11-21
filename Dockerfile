FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pr-api ./cmd/api

FROM alpine:3.19

WORKDIR /root/
COPY --from=builder /app/pr-api .

EXPOSE 4000
CMD ["./pr-api"]
