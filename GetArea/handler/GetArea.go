package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
	"github.com/micro/go-micro/util/log"
	GETAREA "sss/GetArea/proto/GetArea"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"time"
)

type GetArea struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetArea) Call(ctx context.Context, req *GETAREA.Request, rsp *GETAREA.Response) error {
	beego.Info("请求地区信息 /api/v1/areas")
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Error)

	// 1. 从缓存中获取数据， 如果有数据直接发送给前端
	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败: ", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
	}

	// 制定一个key用于在缓存中存取
	area_value := bm.Get("area_info")
	if area_value != nil {
		beego.Info("获取到地区缓存信息: ", string(area_value.([]byte)))
		area_map := []map[string]interface{}{}
		json.Unmarshal(area_value.([]byte), &area_map)

		for key, value := range area_map {
			beego.Info(key, value, "--")
			tmp := GETAREA.Response_Areas{
				Aid:   int32(value["aid"].(float64)),
				Aname: value["aname"].(string),
			}
			rsp.Data = append(rsp.Data, &tmp)
		}
		return nil
	}

	// 2。 缓存中没有数据, 从mysql中查询areas数据
	// beego操作数据库的orm方法
	// 创建orm句柄
	o := orm.NewOrm()
	qs := o.QueryTable("t_area")
	var areas []models.Area
	nums, err := qs.All(&areas)
	if err != nil {
		beego.Info("数据库查询失败")
		rsp.Error = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	if nums == 0 {
		beego.Info("数据库没有数据")
		rsp.Error = utils.RECODE_NODATA
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	// 3。 将数据打包成json字符串存入缓存
	area_json, _ := json.Marshal(areas)
	err = bm.Put("area_info", area_json, time.Second*3600)
	if err != nil {
		beego.Info("数据缓存失败：", err)
		rsp.Error = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)

	}

	// 4。 返回低于信息json给前段
	// 将查询到的数据按照proto的格式发送给web服务
	for key, value := range areas {
		beego.Info("== ", key, value)
		tmp := GETAREA.Response_Areas{
			Aid:   int32(value.Id),
			Aname: value.Name,
		}
		rsp.Data = append(rsp.Data, &tmp)
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetArea) Stream(ctx context.Context, req *GETAREA.StreamingRequest, stream GETAREA.GetArea_StreamStream) error {
	log.Logf("Received GetArea.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GETAREA.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetArea) PingPong(ctx context.Context, stream GETAREA.GetArea_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GETAREA.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
