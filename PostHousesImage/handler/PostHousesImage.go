package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"path"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"

	PostHousesImage_ "sss/PostHousesImage/proto/PostHousesImage"
)

type PostHousesImage struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostHousesImage) Call(ctx context.Context, req *PostHousesImage_.Request, rsp *PostHousesImage_.Response) error {
	beego.Info("用户上传的房屋图片 url: /api/v1.0/houses/:id/images")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	filepath, err := utils.UploadByBuffer(req.Image, path.Ext(req.FileName)[1:])
	if err != nil {
		beego.Info("Postupavatar  models.UploadByBuffer err: ", err)
		rsp.Errno = utils.RECODE_IOERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	sessionId := req.SessionID + "userID"
	valueId := bm.Get(sessionId)
	id, _ := redis.Int(valueId, nil)
	user := models.User{Id: id}

	houseId, _ := strconv.Atoi(req.HouseId)
	house := models.House{Id: houseId, User: &user}

	o := orm.NewOrm()
	err = o.Read(&house, "id")
	if err != nil {
		beego.Info("数据库查询失败：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 	判断index_image url 是否为空
	if house.Index_images_url == "" {
		house.Index_images_url = filepath
	}

	houseImage := models.HouseImage{House: &house, Url: filepath}
	house.Images = append(house.Images, &houseImage)

	_, err = o.Insert(&houseImage)
	if err != nil {
		beego.Info("数据插入失败：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	_, err = o.Update(&house)
	if err != nil {
		beego.Info("房屋数据更新失败：", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.ImageUrl = filepath
	return nil
}
