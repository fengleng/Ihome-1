package main

import (
	"fmt"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"net/http"
	"sss/IhomeWeb/handler"
	_ "sss/IhomeWeb/models"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8999"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("===-------------------------")
	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))
	// 获取地区信息
	service.HandleFunc("/api/v1.0/areas", handler.GetArea)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
