FROM golang:1.23 AS builder

WORKDIR /receiver
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /receiver/main ./cmd

FROM debian:bookworm-slim
WORKDIR /receiver
COPY .env ./cmd/.env
COPY --from=builder /receiver/main ./main
EXPOSE 8000
CMD ["./main"]