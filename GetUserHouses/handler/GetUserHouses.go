package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	GetUserHouses_ "sss/GetUserHouses/proto/GetUserHouses"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
)

type GetUserHouses struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserHouses) Call(ctx context.Context, req *GetUserHouses_.Request, rsp *GetUserHouses_.Response) error {
	beego.Info("获取用户发布的房屋信息 url: /api/v1.0/user/houses")
	//初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 获取sessionid
	sessionId := req.SessionID + "userID"
	valueId := bm.Get(sessionId)
	id, _ := redis.Int(valueId, nil)

	o := orm.NewOrm()
	qs := o.QueryTable("t_house")
	houseList := []models.House{}
	_, err = qs.Filter("user_id", id).All(&houseList)
	if err != nil {
		beego.Info("数据库查询失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// json 编译成二进制数据
	houses, _ := json.Marshal(houseList)
	// 返回二进制数据
	rsp.Mix = houses

	return nil
}
