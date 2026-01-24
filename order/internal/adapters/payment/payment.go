package payment_adapter

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/itallume/microservices-proto/golang/payment"
	"github.com/itallume/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
			retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			retry.WithMax(5),
			retry.WithBackoff(retry.BackoffLinear(time.Second)),
		)))
	conn, err := grpc.Dial(paymentServiceUrl, opts...)

	if err != nil {
		return nil, err
	}

	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := a.payment.Create(ctx,
		&payment.CreatePaymentRequest{
			UserId:     order.CustomerID,
			OrderId:    order.ID,
			TotalPrice: order.TotalPrice(),
		})

	code := status.Code(err)
	if code == codes.DeadlineExceeded {
		log.Printf("2-second deadline exceeded")
		return status.New(codes.DeadlineExceeded, "2-second deadline exceeded").Err()
	}
	return err
}
