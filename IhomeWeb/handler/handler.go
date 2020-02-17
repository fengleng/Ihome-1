package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"image"
	"image/png"
	"net/http"
	go_micro_srv_GetArea "sss/GetArea/proto/GetArea"
	go_micro_srv_GetImageCd "sss/GetImageCd/proto/GetImageCd"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
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

func GetIndex(w http.ResponseWriter, r *http.Request) {
	beego.Info("获取首页轮播图 url:/api/v1.0/house/index")
	response := map[string]interface{}{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetSession(w http.ResponseWriter, r *http.Request) {
	beego.Info("获取Session url:/api/v1.0/session")
	response := map[string]interface{}{
		"errno":  utils.RECODE_SESSIONERR,
		"errmsg": utils.RecodeText(utils.RECODE_SESSIONERR),
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetImages(w http.ResponseWriter, r *http.Request) {
	beego.Info("获取首页轮播图 url: /api/v1.0/imagecode/:uuid")
	uuid := r.URL.Query()["uuid"][0]
	fmt.Println("uuid: == ", uuid)

	server := grpc.NewService()
	server.Init()

	client := go_micro_srv_GetImageCd.NewGetImageCdService("go.micro.srv.GetImageCd", server.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_GetImageCd.Request{
		Uuid: uuid,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var img image.RGBA
	for _, value := range resp.Pix {
		img.Pix = append(img.Pix, uint8(value))
	}
	img.Stride = int(resp.Stride)
	img.Rect.Min.X = int(resp.Min.X)
	img.Rect.Min.Y = int(resp.Min.Y)
	img.Rect.Max.X = int(resp.Max.X)
	img.Rect.Max.Y = int(resp.Max.Y)

	var image captcha.Image
	image.RGBA = &img
	fmt.Println(image)

	png.Encode(w, image)

}
