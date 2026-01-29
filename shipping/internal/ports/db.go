package ports

import "github.com/itallume/microservices/shipping/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Shipping, error)
	Save(shipping *domain.Shipping) error
}