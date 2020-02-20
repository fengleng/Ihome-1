package handler

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"

	GetHouses_ "sss/GetHouses/proto/GetHouses"
)

type GetHouses struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetHouses) Call(ctx context.Context, req *GetHouses_.Request, rsp *GetHouses_.Response) error {
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	areaId, _ := strconv.Atoi(req.AreaId)
	page, _ := strconv.Atoi(req.Page)

	houses := []models.House{}
	o := orm.NewOrm()
	qs := o.QueryTable("t_house")
	num, err := qs.Filter("area_id", areaId).All(&houses)
	if err != nil {
		rsp.Errno = utils.RECODE_PARAMERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	totalPage := int(num)/models.HOUSE_LIST_PAGE_CAPACITY + 1
	housePage := page

	houseList := []interface{}{}
	for _, house := range houses {
		o.LoadRelated(&house, "Area")
		o.LoadRelated(&house, "User")
		o.LoadRelated(&house, "Images")
		o.LoadRelated(&house, "Facilities")
		house.Index_images_url = utils.AddDomain2Url(house.Index_images_url)
		house.User.Avatar_url = utils.AddDomain2Url(house.User.Avatar_url)
		houseList = append(houseList, house)
	}

	rsp.TotalPage = int64(totalPage)
	rsp.CurrentPage = int64(housePage)
	rsp.Houses, _ = json.Marshal(&houseList)

	return nil
}
