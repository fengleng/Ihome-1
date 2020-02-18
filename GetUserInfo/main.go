package main

import (
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro"
	"sss/GetUserInfo/handler"
	"sss/GetUserInfo/subscriber"

	GetUserInfo "sss/GetUserInfo/proto/GetUserInfo"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.GetUserInfo"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	GetUserInfo.RegisterGetUserInfoHandler(service.Server(), new(handler.GetUserInfo))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetUserInfo", service.Server(), new(subscriber.GetUserInfo))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetUserInfo", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
