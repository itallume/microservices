package shipping_adapter

import (
	"context"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/itallume/microservices-proto/golang/shipping"
	"github.com/itallume/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(
			retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			retry.WithMax(5),
			retry.WithBackoff(retry.BackoffLinear(time.Second)),
		)))
	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) CalculateRoute(order *domain.Order, billId int64) (int32, error) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	var items []*shipping.Item
	for _, orderItem := range order.OrderItems {
		items = append(items, &shipping.Item{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	response, err := a.shipping.Create(ctx,
		&shipping.CreateShippingRequest{
			BillId:     billId,
			OrderItems: items,
		})

	code := status.Code(err)
	if code == codes.DeadlineExceeded {
		log.Printf("2-second deadline exceeded")
		return 0, status.New(codes.DeadlineExceeded, "2-second deadline exceeded").Err()
	}
	return response.DeliveryTime, err
}
