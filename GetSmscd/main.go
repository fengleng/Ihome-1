package main

import (
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"sss/GetSmscd/handler"
	"sss/GetSmscd/subscriber"

	GetSmscd "sss/GetSmscd/proto/GetSmscd"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetSmscd"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetSmscd.RegisterGetSmscdHandler(service.Server(), new(handler.GetSmscd))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSmscd", service.Server(), new(subscriber.GetSmscd))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSmscd", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
