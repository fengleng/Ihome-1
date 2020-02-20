package main

import (
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"sss/PostHouses/handler"
	"sss/PostHouses/subscriber"

	PostHouses "sss/PostHouses/proto/PostHouses"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostHouses"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	PostHouses.RegisterPostHousesHandler(service.Server(), new(handler.PostHouse))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostHouses", service.Server(), new(subscriber.PostHouse))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostHouses", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
