package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"
	"time"

	GetHouseInfo_ "sss/GetHouseInfo/proto/GetHouseInfo"
)

type GetHouseInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetHouseInfo) Call(ctx context.Context, req *GetHouseInfo_.Request, rsp *GetHouseInfo_.Response) error {
	beego.Info("查看房屋详情 url: /api/v1.0/houses/:id")

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	sessionId := req.SessionId + "userID"
	valueId := bm.Get(sessionId)
	id, _ := redis.Int(valueId, nil)
	//user := models.User{Id: id}

	houseId, _ := strconv.Atoi(req.HouseId)
	houseInfoKey := fmt.Sprintf("house_info_%s", houseId)
	houseIfoValue := bm.Get(houseInfoKey)
	if houseIfoValue != nil {
		beego.Info("房屋信息在缓存中")
		rsp.UserId = int64(id)
		rsp.HouseData = houseIfoValue.([]byte)
	}

	o := orm.NewOrm()
	house := models.House{Id: houseId}
	o.Read(&house, "id")

	// 关联查询
	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Facilities")

	// 将查询的结果存入redis中
	houseMix, err := json.Marshal(house)
	bm.Put(houseInfoKey, houseMix, time.Second*3600)

	rsp.UserId = int64(id)
	rsp.HouseData = houseMix

	return nil
}
