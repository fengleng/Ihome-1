package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"time"

	GetIndex_ "sss/GetIndex/proto/GetIndex"
)

type GetIndex struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetIndex) Call(ctx context.Context, req *GetIndex_.Request, rsp *GetIndex_.Response) error {
	beego.Info("获取首页轮播图 url:/api/v1.0/house/index")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	index_page_key := "home_page_data"
	value := bm.Get(index_page_key)
	if value != nil {
		beego.Info("index_data_key 在redis缓存中")
		rsp.Mix = value.([]byte)
		return nil
	}

	o := orm.NewOrm()
	houses := []models.House{}
	if _, err := o.QueryTable("t_house").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses); err != nil {
		beego.Info("数据库查询失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	data := []interface{}{}
	for _, house := range houses {
		o.LoadRelated(&house, "Area")
		o.LoadRelated(&house, "User")
		o.LoadRelated(&house, "Images")
		o.LoadRelated(&house, "Facilities")
		data = append(data, house.To_housr_info())
	}

	indexPageValue, _ := json.Marshal(data)
	bm.Put(index_page_key, indexPageValue, time.Second*3600)
	rsp.Mix = indexPageValue
	return nil
}
