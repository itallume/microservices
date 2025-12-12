package ports

import "github.com/itallume/microservices/order/internal/application/core/domain"

type APIPort interface{
	PlaceOrder(order domain.Order) (domain.Order, error)
}