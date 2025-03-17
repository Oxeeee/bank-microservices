FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

FROM debian:buster-slim

WORKDIR /app

COPY --from=builder /app/app .
COPY configs/dev.yaml /app/configs/dev.yaml
COPY templates/ /app/templates

ENV CONFIG_PATH=/app/configs/dev.yaml

EXPOSE 8080

CMD ["./app"]