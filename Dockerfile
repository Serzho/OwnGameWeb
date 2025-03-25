FROM golang:1.24 AS builder

COPY . .
COPY vendor ./vendor
RUN go build -mod=vendor -o main cmd/app/main.go

ENTRYPOINT ["./main"]
# CMD tail -f /dev/null