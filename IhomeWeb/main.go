package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
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

	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))
	// 获取地区信息
	rou.GET("/api/v1.0/areas", handler.GetArea)
	// 获取首页轮播图
	rou.GET("/api/v1.0/house/index", handler.GetIndex)
	// 获取session
	rou.GET("/api/v1.0/session", handler.GetSession)
	// 获取验证码图片
	rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImages)
	// 获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile", handler.GetSmsCode)
	// 用户注册
	rou.POST("/api/v1.0/users", handler.PostRet)
	// 用户登录
	rou.POST("/api/v1.0/sessions", handler.PostLogin)
	// 用户退出
	rou.DELETE("/api/v1.0/session", handler.DeleteLogout)
	// 用户个人信息查看
	rou.GET("/api/v1.0/user", handler.GetUserInfo)
	// 上传用户头像
	rou.POST("/api/v1.0/user/avatar", handler.PostAvatar)
	// 更新用户信息
	rou.PUT("/api/v1.0/user/name", handler.PutUserInfo)
	// 检查用户实名认证
	rou.GET("/api/v1.0/user/auth", handler.GetUserInfo)
	// 更新实名认证信息
	rou.POST("/api/v1.0/user/auth", handler.PostUserAuth)
	// 获取用户发布的房源
	rou.GET("/api/v1.0/user/houses", handler.GetUserHouses)
	// 发布新房源
	rou.POST("/api/v1.0/houses", handler.PostHouses)
	// 房屋图片上传
	rou.POST("/api/v1.0/houses/:id/images", handler.PostHousesImage)
	// 查看房屋详细信息
	rou.GET("/api/v1.0/houses/:id", handler.GetHouseInfo)


	service.Handle("/", rou)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
