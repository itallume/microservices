package api

import (
	"github.com/itallume/microservices/order/internal/application/core/domain"
	"github.com/itallume/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	quantity := 0
	for _, item := range order.OrderItems {
		quantity += int(item.Quantity)
		if quantity > 50 {
			return domain.Order{}, status.Errorf(codes.InvalidArgument, "Orders must have a maximum of 50 items.")
		}
	}

	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}
	return order, nil
}
