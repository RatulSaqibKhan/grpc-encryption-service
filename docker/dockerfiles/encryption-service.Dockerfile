FROM golang:1.23-alpine

WORKDIR /app

COPY ./apps/encryption-service/go.mod ./apps/encryption-service/go.sum ./
RUN go mod download

COPY ./apps/encryption-service .

RUN go build -o grpc-service ./cmd/main.go

CMD ["./grpc-service"]
