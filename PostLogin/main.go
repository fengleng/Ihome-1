package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"sss/PostLogin/handler"
	"sss/PostLogin/subscriber"

	PostLogin "sss/PostLogin/proto/PostLogin"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostLogin"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	PostLogin.RegisterPostLoginHandler(service.Server(), new(handler.PostLogin))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostLogin", service.Server(), new(subscriber.PostLogin))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostLogin", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
