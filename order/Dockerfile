FROM golang:1.24 AS builder
WORKDIR /app

# Copia tudo (proto e microservices)
COPY microservices-proto/ ./microservices-proto/
COPY microservices/order/ ./microservices/order/

WORKDIR /app/microservices/order
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main.go

FROM scratch
COPY --from=builder /app/microservices/order/order ./order
CMD ["./order"]
