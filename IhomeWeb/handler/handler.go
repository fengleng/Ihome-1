package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"net/http"
	go_micro_srv_GetArea "sss/GetArea/proto/GetArea"
	"sss/IhomeWeb/models"
)

func GetArea(w http.ResponseWriter, r *http.Request) {
	beego.Info("获取地区请求客户端 url:api/v1.0/areas")
	server := grpc.NewService(micro.Name("go.micro.web.IhomeWeb"))
	server.Init()

	// 调用服务返回句柄
	client := go_micro_srv_GetArea.NewGetAreaService("go.micro.srv.GetArea", server.Client())

	resp, err := client.Call(context.TODO(), &go_micro_srv_GetArea.Request{})
	if err != nil {
		beego.Info("err: == ", err)
		beego.Info("resp: == ", resp)
		http.Error(w, err.Error(), 500)
		return
	}

	area_list := []models.Area{}

	for _, value := range resp.Data {
		tmp := models.Area{
			Id:   int(value.Aid),
			Name: value.Aname,
		}
		area_list = append(area_list, tmp)
	}

	response := map[string]interface{}{
		"errno":  resp.Error,
		"errmsg": resp.Errmsg,
		"data":   area_list,
	}

	// 回传数据的时候需要设置数据格式
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}
