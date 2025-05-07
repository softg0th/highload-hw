FROM golang:1.23 AS builder

WORKDIR /storage
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /storage/main ./cmd


FROM debian:bookworm-slim
WORKDIR /storage
COPY .env ./cmd/.env
COPY --from=builder /storage/main ./main
EXPOSE 8002
EXPOSE 50051
CMD ["./main"]