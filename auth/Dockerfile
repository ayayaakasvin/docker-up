FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o auth-binary /app/cmd/auth/main.go

FROM alpine:latest

RUN apk --no-cache add libc6-compat

COPY --from=builder /app/auth-binary .

COPY --from=builder /app/config/config.yaml /app/config/config.yaml

RUN chmod +x ./auth-binary

EXPOSE 8088

CMD ["./auth-binary"]