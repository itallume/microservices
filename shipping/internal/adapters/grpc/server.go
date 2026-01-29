package grpc

import (
	"context"
	"github.com/itallume/microservices/shipping/internal/ports"
	"github.com/itallume/microservices-proto/golang/shipping"
	"github.com/itallume/microservices/shipping/internal/application/core/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"fmt"
	"net"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/itallume/microservices/shipping/config"
)


type Adapter struct {
	api ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter{
	return &Adapter{api: api, port: port}
}

func (a Adapter) Create(ctx context.Context, request *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error){
	var items []domain.Item

	for _, item := range request.OrderItems{
		items = append(items, domain.Item{
			ProductCode: item.ProductCode,
			UnitPrice: item.UnitPrice,
			Quantity: item.Quantity,
		})
	}

	newShipping := domain.NewShipping(request.BillId, items)
	result, err := a.api.CalculateRoute(newShipping)

	code := status.Code(err)
	if code == codes.InvalidArgument {
		return nil, err
	} else if err != nil{
		return nil, status.New(codes.Internal, fmt.Sprintf("Failed to calculate route. %v ", err)).Err()
	}

	deliveryTime, err := result.CalculateDeliveryTime()
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("Failed to calculate delivery time. %v", err)).Err()
	}

	return &shipping.CreateShippingResponse{DeliveryTime: deliveryTime}, nil
}

func (a Adapter) Run(){
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil{
		log.Fatalf("failed to listen on port %d, error : %v", a.port , err)
	}	
	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)
	if config.GetEnv() == "development"{
		reflection.Register(grpcServer)
	}
	if err := grpcServer.Serve(listen); err != nil{
		log.Fatalf ("failed to serve grpc on port")
	}
}
