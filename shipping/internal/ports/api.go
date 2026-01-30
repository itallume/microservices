package ports

import "github.com/itallume/microservices/shipping/internal/application/core/domain"

type APIPort interface{
	CalculateRoute(shipping domain.Shipping, billId int64) (int32, error)
}