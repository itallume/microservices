package ports

import "github.com/itallume/microservices/shipping/internal/application/core/domain"

type APIPort interface{
	CalculateRoute(shipping domain.Shipping) (domain.Shipping, error)
}