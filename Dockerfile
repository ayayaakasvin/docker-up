FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/restapi

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -o backend ./cmd/restapi/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add libc6-compat

COPY --from=builder /app/restapi/backend .

COPY --from=builder /app/restapi/config/config.yaml /app/restapi/config/config.yaml

COPY --from=builder /app/restapi/migration /app/restapi/migration

RUN chmod +x ./backend

CMD ["./backend"]