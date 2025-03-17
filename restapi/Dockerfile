FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o backend /app/cmd/restapi/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add libc6-compat

COPY --from=builder /app/backend .

COPY --from=builder /app/config/config.yaml /app/config/config.yaml

RUN chmod +x ./backend

EXPOSE 8069

CMD ["./backend"]