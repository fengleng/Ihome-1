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
		web.Address(":8008"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("===-------------------------")

	//rou := httprouter.New()
	//rou.NotFound = http.FileServer(http.Dir("html"))
	//rou.GET("/api/v1.0/areas", handler.GetArea)
	//rou.GET("/api/v1.0/house/index", handler.GetIndex)
	//rou.GET("/api/v1.0/session", handler.GetSession)

	// register html handler
	service.Handle("/", http.FileServer(http.Dir("html")))
	// 获取地区信息
	service.HandleFunc("/api/v1.0/areas", handler.GetArea)
	// 获取首页轮播图
	service.HandleFunc("/api/v1.0/house/index", handler.GetIndex)
	// 获取session
	service.HandleFunc("/api/v1.0/session", handler.GetSession)

	// 获取验证码图片
	service.HandleFunc("/api/v1.0/imagecode", handler.GetImages)


	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
