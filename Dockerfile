# Build container
FROM golang:1.24.2-bookworm AS builder
WORKDIR /app

COPY . .
RUN go mod download ; \
    go build -o main ./cmd/EmailFilterClient

# Final container
FROM debian:bookworm-slim
WORKDIR /app

RUN apt update ; \
    apt install ca-certificates -y

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

ENV PORT=8080
ENV BASIC_AUTH_PASSWORD=mailAdminPwd1!

CMD ["sh", "-c", "/app/main -port=$PORT -basicAuthPassword=$BASIC_AUTH_PASSWORD"]

EXPOSE 8080
