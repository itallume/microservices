package api

import (

	"github.com/itallume/microservices/shipping/internal/application/core/domain"
	"github.com/itallume/microservices/shipping/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) CalculateRoute(shipping domain.Shipping) (domain.Shipping, error) {
	if shipping.ItemsQuantity() <= 0 {
		return domain.Shipping{}, status.Errorf(codes.InvalidArgument, "Shipping must have at least one item.")
	}
	err := a.db.Save(&shipping)
	if err != nil {
		return domain.Shipping{}, err
	}
	return shipping, nil
}