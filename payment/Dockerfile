FROM golang:1.24 AS builder
WORKDIR /app

# Copia tudo (proto e microservices)
COPY microservices-proto/ ./microservices-proto/
COPY microservices/payment/ ./microservices/payment/

WORKDIR /app/microservices/payment
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go

FROM scratch
COPY --from=builder /app/microservices/payment/payment ./payment
CMD ["./payment"]