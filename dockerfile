FROM golang:1.24.1-alpine3.21 AS builder
WORKDIR /app
COPY ["./", "./"]
RUN go build -o ./bin cmd/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/bin .
EXPOSE 8080

CMD ["./bin"]