FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mini-redis ./main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/mini-redis .

EXPOSE 6379

ENTRYPOINT ["./mini-redis"]
