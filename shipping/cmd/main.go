package main

import (
	"log"

	"github.com/itallume/microservices/shipping/config"
	"github.com/itallume/microservices/shipping/internal/adapters/db"

	//"github.com/itallume/microservices/order/internal/adapters/rest"
	"github.com/itallume/microservices/shipping/internal/adapters/grpc"
	"github.com/itallume/microservices/shipping/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}