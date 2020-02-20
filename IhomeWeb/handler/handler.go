package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"github.com/julienschmidt/httprouter"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/service/grpc"
	"image"
	"image/png"
	"net/http"
	"regexp"
	go_micro_srv_DeleteSession "sss/DeleteSession/proto/DeleteSession"
	go_micro_srv_GetArea "sss/GetArea/proto/GetArea"
	go_micro_srv_GetImageCd "sss/GetImageCd/proto/GetImageCd"
	go_micro_srv_GetSession "sss/GetSession/proto/GetSession"
	go_micro_srv_GetSmscd "sss/GetSmscd/proto/GetSmscd"
	go_micro_srv_GetUserHouses "sss/GetUserHouses/proto/GetUserHouses"
	go_micro_srv_GetUserInfo "sss/GetUserInfo/proto/GetUserInfo"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	go_micro_srv_PostAvatar "sss/PostAvatar/proto/PostAvatar"
	go_micro_srv_PostHouses "sss/PostHouses/proto/PostHouses"
	go_micro_srv_PostHousesImage "sss/PostHousesImage/proto/PostHousesImage"
	go_micro_srv_PostLogin "sss/PostLogin/proto/PostLogin"
	go_micro_srv_PostRet "sss/PostRet/proto/PostRet"
	go_micro_srv_PostUserAuth "sss/PostUserAuth/proto/PostUserAuth"
	go_micro_srv_PutUserInfo "sss/PutUserInfo/proto/PutUserInfo"
)

type ReqestParam struct {
	Name      string   `json:"name"`
	IDCard    string   `json:"id_card"`
	RealName  string   `json:"real_name"`
	Mobile    string   `json:"mobile"`
	Password  string   `json:"password"`
	SmsCode   string   `json:"sms_code"`
	Title     string   `json:"title"`
	Price     string   `json:"price"`
	AreaId    string   `json:"area_id"`
	Address   string   `json:"address"`
	RoomCount string   `json:"room_count"`
	Acreage   string   `json:"acreage"`
	Unit      string   `json:"unit"`
	Capacity  string   `json:"capacity"`
	Beds      string   `json:"beds"`
	Deposit   string   `json:"deposit"`
	MinDays   string   `json:"min_days"`
	MaxDays   string   `json:"max_days"`
	Facility  []string `json:"facility"`
}

func parseParams(w http.ResponseWriter, r *http.Request) (*ReqestParam, error) {
	var reqParams = &ReqestParam{}
	err := json.NewDecoder(r.Body).Decode(reqParams)
	return reqParams, err
}

func initService() micro.Service {
	service := grpc.NewService()
	service.Init()
	return service
}

func GetArea(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	beego.Info("获取地区请求客户端 url:api/v1.0/areas")
	service := initService()

	// 调用服务返回句柄
	client := go_micro_srv_GetArea.NewGetAreaService("go.micro.srv.GetArea", service.Client())

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

func GetIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func GetSession(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取Session url:/api/v1.0/session")
	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	service := initService()

	client := go_micro_srv_GetSession.NewGetSessionService("go.micro.srv.GetSession", service.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_GetSession.Request{
		SessionID: cookie.Value,
	})
	if err != nil {
		utils.Response(w, resp.Errno, utils.RecodeText(resp.Errmsg), nil)
		return
	}
	data := map[string]string{}
	data["name"] = resp.UserName
	beego.Info("========================")
	if err := utils.Response(w, resp.Errno, utils.RecodeText(resp.Errmsg), data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	return
}

func GetImages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取首页轮播图 url: /api/v1.0/imagecode/:uuid")
	uuid := p.ByName("uuid")
	fmt.Println("uuid: == ", uuid)

	service := initService()

	client := go_micro_srv_GetImageCd.NewGetImageCdService("go.micro.srv.GetImageCd", service.Client())
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

func GetSmsCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取短信验证码 url: /api/v1.0/smscode/:mobile")
	mobile := p.ByName("mobile")
	// 正则表达式 手机号
	mobileReg := regexp.MustCompile(`0?(13|14|15|17|18|19)[0-9]{9}`)
	bl := mobileReg.MatchString(mobile)
	if !bl {
		beego.Info("+++++", bl)
		utils.Response(w, utils.RECODE_MOBILEERR, utils.RecodeText(utils.RECODE_MOBILEERR), nil)
		return

	}
	beego.Info("====================")

	imageStr := r.URL.Query()["text"][0]
	uuid := r.URL.Query()["id"][0]

	service := initService()

	client := go_micro_srv_GetSmscd.NewGetSmscdService("go.micro.srv.GetSmscd", service.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_GetSmscd.Request{
		Mobile:   mobile,
		Uuid:     uuid,
		ImageStr: imageStr,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := utils.Response(w, resp.Error, resp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func PostRet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户注册 url: api/v1.0/users")
	reqParams, _ := parseParams(w, r)
	beego.Info(reqParams, "===========")

	if reqParams.Mobile == "" || reqParams.Password == "" || reqParams.SmsCode == "" {
		utils.Response(w, utils.RECODE_PARAMERR, utils.RecodeText(utils.RECODE_PARAMERR), nil)
		return
	}

	service := initService()

	client := go_micro_srv_PostRet.NewPostRetService("go.micro.srv.PostRet", service.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_PostRet.Request{
		Mobile:   reqParams.Mobile,
		Password: reqParams.Password,
		SmsCode:  reqParams.SmsCode,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// 读取cookie 统一 "userLogin"
	setCookie(w, r, resp.SessionId)
	if err := utils.Response(w, resp.Erron, resp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户登陆 url: /api/v1.0/sessions")
	reqParams, _ := parseParams(w, r)
	beego.Info("请求数据：", reqParams)
	service := initService()
	client := go_micro_srv_PostLogin.NewPostLoginService("go.micro.srv.PostLogin", service.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_PostLogin.Request{
		Mobile:   reqParams.Mobile,
		Password: reqParams.Password,
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	setCookie(w, r, resp.SessionID)
	if err := utils.Response(w, resp.Errno, resp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func setCookie(w http.ResponseWriter, r *http.Request, sessionId string) {
	// 读取cookie 统一 "userLogin"
	cookie, err := r.Cookie("userLogin")
	if err != nil || "" == cookie.Value {
		cookie := &http.Cookie{
			Name:   "userLogin",
			Value:  sessionId,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	}
}

func DeleteLogout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户登陆 url: /api/v1.0/session")
	userlogin, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	service := initService()
	client := go_micro_srv_DeleteSession.NewDeleteSessionService("go.micro.srv.DeleteSession", service.Client())
	resp, err := client.Call(context.TODO(), &go_micro_srv_DeleteSession.Request{
		SessionID: userlogin.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	cookie, err := r.Cookie("userLogin")
	if err != nil || cookie.Value == "" {
		return
	} else {
		cookie := http.Cookie{Name: "userLogin", Path: "/", MaxAge: -1}
		http.SetCookie(w, &cookie)
	}

	if err := utils.Response(w, resp.Errno, resp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 503)
		beego.Info(err)
		return
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户信息 url: /api/v1.0/user")
	service := initService()

	client := go_micro_srv_GetUserInfo.NewGetUserInfoService("go.micro.srv.GetUserInfo", service.Client())
	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	rsp, err := client.Call(context.TODO(), &go_micro_srv_GetUserInfo.Request{
		SessionID: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	data := map[string]interface{}{
		"user_id":    rsp.UserID,
		"name":       rsp.Name,
		"mobile":     rsp.Mobile,
		"real_name":  rsp.RealName,
		"id_card":    rsp.IDCard,
		"avatar_url": utils.AddDomain2Url(rsp.AvatarURL),
	}
	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, data); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

func PostAvatar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("上传用户头像 url: /api/v1.0/user/avatar")

	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	file, handle, err := r.FormFile("avatar")
	if err != nil {
		utils.Response(w, utils.RECODE_IOERR, utils.RecodeText(utils.RECODE_IOERR), nil)
		return
	}

	fileBuffer := make([]byte, handle.Size)

	_, err = file.Read(fileBuffer)
	if err != nil {
		utils.Response(w, utils.RECODE_IOERR, utils.RecodeText(utils.RECODE_IOERR), nil)
		return
	}

	service := initService()
	client := go_micro_srv_PostAvatar.NewPostAvatarService("go.micro.srv.PostAvatar", service.Client())
	rsp, err := client.Call(context.TODO(), &go_micro_srv_PostAvatar.Request{
		Avatar:    fileBuffer,
		SessionID: cookie.Value,
		FileSize:  handle.Size,
		FileName:  handle.Filename,
	})

	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	data := make(map[string]string)
	data["avatar_url"] = utils.AddDomain2Url(rsp.AvatarUrl)

	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, data); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

}

func PutUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("更新用户姓名 url: /api/v1.0/user/name")

	reqParams, _ := parseParams(w, r)

	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	service := initService()
	client := go_micro_srv_PutUserInfo.NewPutUserInfoService("go.micro.srv.PutUserInfo", service.Client())
	rsp, err := client.Call(context.TODO(), &go_micro_srv_PutUserInfo.Request{
		SessionID: cookie.Value,
		UserName:  reqParams.Name,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	data := map[string]string{
		"name": rsp.UserName,
	}

	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, data); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

}

func PostUserAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("更新实名认证 url: /api/v1.0/user/auth")
	reqParams, _ := parseParams(w, r)

	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	service := initService()
	client := go_micro_srv_PostUserAuth.NewPostUserAuthService("go.micro.srv.PostUserAuth", service.Client())
	rsp, err := client.Call(context.TODO(), &go_micro_srv_PostUserAuth.Request{
		RealName:  reqParams.RealName,
		IDCard:    reqParams.IDCard,
		SessionID: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, nil); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

}

func GetUserHouses(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("获取用户发布的房屋信息 url: /api/v1.0/user/houses")

	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	service := initService()
	client := go_micro_srv_GetUserHouses.NewGetUserHousesService("go.micro.srv.GetUserHouses", service.Client())
	rsp, err := client.Call(context.TODO(), &go_micro_srv_GetUserHouses.Request{
		SessionID: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	housesList := []models.House{}
	err = json.Unmarshal(rsp.Mix, &housesList)
	if err != nil {
		beego.Info("反序列化异常：", err)

	}
	var houses []interface{}
	for _, value := range housesList {
		houses = append(houses, value.To_housr_info())
	}

	data := map[string]interface{}{
		"houses": houses,
	}
	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, data); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

}

func PostHouses(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户发布的房屋信息 url: /api/v1.0/user/houses")

	reqParams := ReqestParam{}
	json.NewDecoder(r.Body).Decode(&reqParams)

	fmt.Printf("===========, %+v\n", reqParams)
	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	service := initService()

	client := go_micro_srv_PostHouses.NewPostHousesService("go.micro.srv.PostHouses", service.Client())
	rsp, err := client.Call(context.TODO(), &go_micro_srv_PostHouses.Request{
		Title:     reqParams.Title,
		Price:     reqParams.Price,
		AreaId:    reqParams.AreaId,
		Address:   reqParams.Address,
		RoomCount: reqParams.RoomCount,
		Acreage:   reqParams.Acreage,
		Unit:      reqParams.Unit,
		Capacity:  reqParams.Capacity,
		Beds:      reqParams.Beds,
		Deposit:   reqParams.Deposit,
		MinDays:   reqParams.MinDays,
		MaxDays:   reqParams.MaxDays,
		Facility:  reqParams.Facility,
		SessionId: cookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	data := map[string]interface{}{
		"house_id": rsp.HouseId,
	}
	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, data); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}
}

func PostHousesImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	beego.Info("用户上传的房屋图片 url: /api/v1.0/houses/:id/images")
	houseId := p.ByName("id")
	cookie, err := r.Cookie("userLogin")
	if err != nil {
		utils.Response(w, utils.RECODE_SESSIONERR, utils.RecodeText(utils.RECODE_SESSIONERR), nil)
		return
	}

	file, handle, err := r.FormFile("house_image")
	if err != nil {
		utils.Response(w, utils.RECODE_IOERR, utils.RecodeText(utils.RECODE_IOERR), nil)
		return
	}

	fileBuffer := make([]byte, handle.Size)
	_, err = file.Read(fileBuffer)
	if err != nil {
		utils.Response(w, utils.RECODE_IOERR, utils.RecodeText(utils.RECODE_IOERR), nil)
		return
	}

	service := initService()
	client := go_micro_srv_PostHousesImage.NewPostHousesImageService("go.micro.srv.PostHousesImage", service.Client())
	rsp, err := client.Call(context.TODO(), &go_micro_srv_PostHousesImage.Request{
		Image:     fileBuffer,
		SessionID: cookie.Value,
		FileSize:  handle.Size,
		FileName:  handle.Filename,
		HouseId:   houseId,
	})
	if err != nil {
		http.Error(w, err.Error(), 502)
		return
	}

	data := map[string]string{
		"url": utils.AddDomain2Url(rsp.ImageUrl),
	}

	if err := utils.Response(w, rsp.Errno, rsp.Errmsg, data); err != nil {
		http.Error(w, err.Error(), 503)
		return
	}

}
