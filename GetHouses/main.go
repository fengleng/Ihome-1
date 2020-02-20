package main

import (
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"sss/GetHouses/handler"
	"sss/GetHouses/subscriber"

	GetHouses "sss/GetHouses/proto/GetHouses"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetHouses"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetHouses.RegisterGetHousesHandler(service.Server(), new(handler.GetHouses))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetHouses", service.Server(), new(subscriber.GetHouses))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetHouses", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
