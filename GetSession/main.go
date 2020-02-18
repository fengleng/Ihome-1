package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"sss/GetSession/handler"
	"sss/GetSession/subscriber"

	GetSession "sss/GetSession/proto/GetSession"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetSession"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetSession.RegisterGetSessionHandler(service.Server(), new(handler.GetSession))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSession", service.Server(), new(subscriber.GetSession))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetSession", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
