package ports

import "github.com/itallume/microservices/order/internal/application/core/domain"

type ShippingPort interface{
	CalculateRoute(shipping *domain.Order, billId int64) (int32, error)
}