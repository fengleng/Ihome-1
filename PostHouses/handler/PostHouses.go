package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"

	PostHouses "sss/PostHouses/proto/PostHouses"
)

type PostHouse struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostHouse) Call(ctx context.Context, req *PostHouses.Request, rsp *PostHouses.Response) error {
	beego.Info("用户发布的房屋信息 url: /api/v1.0/user/houses")
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

	o := orm.NewOrm()

	user := models.User{Id: id}
	areaID, _ := strconv.Atoi(req.AreaId)
	area := models.Area{Id: areaID}
	o.Read(&user, "id")
	o.Read(&area, "id")

	facilities := []*models.Facility{}
	for _, value := range req.Facility {
		fid, _ := strconv.Atoi(value)
		f := &models.Facility{Id: fid}
		facilities = append(facilities, f)
	}

	price, _ := strconv.Atoi(req.Price)
	roomCOunt, _ := strconv.Atoi(req.RoomCount)
	acreage, _ := strconv.Atoi(req.Acreage)
	capacity, _ := strconv.Atoi(req.Capacity)
	deposit, _ := strconv.Atoi(req.Deposit)
	mindays, _ := strconv.Atoi(req.MinDays)
	maxdays, _ := strconv.Atoi(req.MaxDays)
	house := models.House{
		User:       &user,
		Title:      req.Title,
		Price:      price * 100,
		Area:       &area,
		Address:    req.Address,
		Room_count: roomCOunt,
		Acreage:    acreage,
		Unit:       req.Unit,
		Capacity:   capacity,
		Beds:       req.Beds,
		Deposit:    deposit*100,
		Min_days:   mindays,
		Max_days:   maxdays,
		//Facilities: facilities,
	}

	houseId, err := o.Insert(&house)
	if err != nil {
		beego.Info("数据库操作失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 多对多插入
	m2m := o.QueryM2M(&house, "facilities")
	_, err = m2m.Add(facilities)
	if err != nil {
		beego.Info("房屋设施多对多失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.HouseId = strconv.Itoa(int(houseId))

	return nil
}
