package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"sss/GetIndex/handler"
	"sss/GetIndex/subscriber"

	GetIndex "sss/GetIndex/proto/GetIndex"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetIndex"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetIndex.RegisterGetIndexHandler(service.Server(), new(handler.GetIndex))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetIndex", service.Server(), new(subscriber.GetIndex))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetIndex", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
