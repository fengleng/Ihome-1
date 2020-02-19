package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro/util/log"
	"sss/PutUserInfo/handler"
	"sss/PutUserInfo/subscriber"

	PutUserInfo "sss/PutUserInfo/proto/PutUserInfo"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PutUserInfo"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	PutUserInfo.RegisterPutUserInfoHandler(service.Server(), new(handler.PutUserInfo))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PutUserInfo", service.Server(), new(subscriber.PutUserInfo))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PutUserInfo", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
