package api

import (
	"log"

	"github.com/itallume/microservices/order/internal/application/core/domain"
	"github.com/itallume/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
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
	billId, paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}

	deliveryTime, shippingErr := a.shipping.CalculateRoute(&order, billId)
	if shippingErr != nil {
		return domain.Order{}, shippingErr
	}
	log.Printf("Order %d will be delivered in %d days", order.ID, deliveryTime)

	return order, nil
}
