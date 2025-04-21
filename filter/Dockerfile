FROM golang:1.23 AS builder

WORKDIR /filter
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /filter/main ./cmd

FROM debian:bookworm-slim
WORKDIR /filter
COPY .env ./cmd/.env
COPY --from=builder /filter/main ./main
EXPOSE 8000
CMD ["./main"]