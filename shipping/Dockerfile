FROM golang:1.24 AS builder
WORKDIR /app

# Copia tudo (proto e microservices)
COPY microservices-proto/ ./microservices-proto/
COPY microservices/shipping/ ./microservices/shipping/

WORKDIR /app/microservices/shipping
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipping ./cmd/main.go

FROM scratch
COPY --from=builder /app/microservices/shipping/shipping ./shipping
CMD ["./shipping"]
