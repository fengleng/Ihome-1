package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"strconv"
	"strings"
	"time"

	"github.com/micro/go-micro/util/log"

	GetSms "sss/GetSmscd/proto/GetSmscd"
)

type GetSmscd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSmscd) Call(ctx context.Context, req *GetSms.Request, rsp *GetSms.Response) error {
	beego.Info("获取短信验证码 url: /api/v1.0/smscode/:mobile")
	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	// 验证手机号是否已经存在
	o := orm.NewOrm()
	user := models.User{Mobile: req.Mobile}
	err := o.Read(&user, "mobile")
	beego.Info("手机号查询：", err, ", Mobile: ", req.Mobile)
	if err == nil {
		beego.Info("用户已存在")
		rsp.Error = utils.RECODE_MOBILEXISTEERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_MOBILEXISTEERR)
		return nil
	}

	// 判断验证码是否已经存在
	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败: ", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}
	code := bm.Get(req.Uuid)
	codeStr, _ := redis.String(code, nil)
	beego.Info("机器验证码：", strings.ToLower(req.ImageStr), ", 你输入的验证码：", strings.ToLower(codeStr))
	if strings.ToLower(req.ImageStr) != strings.ToLower(codeStr) {
		beego.Info("验证码错误: ")
		rsp.Error = utils.RECODE_VERIFCODEILERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
	}

	// 发送短信
	val := sendSMS(req.Mobile)
	beego.Info("短信验证码：", val)
	if bm.Put(req.Mobile, val, time.Second*300) != nil {
		beego.Info("redis 创建失败: ", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
		return nil
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSmscd) Stream(ctx context.Context, req *GetSms.StreamingRequest, stream GetSms.GetSmscd_StreamStream) error {
	log.Logf("Received GetSmscd.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetSms.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetSmscd) PingPong(ctx context.Context, stream GetSms.GetSmscd_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetSms.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

func sendSMS(mobile string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := strconv.Itoa(r.Intn(9999) + 1001)
	//// SMS 短信服务配置 appid & appkey 请前往：https://www.mysubmail.com/chs/sms/apps 获取
	//config := make(map[string]string)
	//config["appid"] = "46357"
	//config["appkey"] = "04d7eb0a3ad0a718d9bce812174c5353"
	//// SMS 数字签名模式 normal or md5 or sha1 ,normal = 明文appkey鉴权 ，md5 和 sha1 为数字签名鉴权模式
	//config["signType"] = "sha1"
	//
	////创建 短信 Send 接口
	//
	//submail := sms.CreateXsend(config)
	////设置联系人 手机号码
	//submail.SetTo(mobile)
	////设置短信模板id
	//submail.SetProject("7BHce1")
	////添加模板中的设置的动态变量。如模板为：【xxx】您的验证码是:@var(code),请在@var(time)分钟内输入。
	//submail.AddVar("code", code)
	//submail.AddVar("time", "5")
	////执行 Xsend 方法发送短信
	//xsend := submail.Xsend()
	//fmt.Println("短信XSend 接口:", xsend)
	return code
}
