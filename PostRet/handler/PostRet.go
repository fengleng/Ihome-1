package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"time"

	"github.com/micro/go-micro/util/log"

	PostRegiste "sss/PostRet/proto/PostRet"
)

type PostRet struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostRet) Call(ctx context.Context, req *PostRegiste.Request, rsp *PostRegiste.Response) error {
	beego.Info("用户注册 url: api/v1.0/users")
	rsp.Erron = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Erron)

	// 验证短信验证码
	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接错误：", err)
		rsp.Errmsg = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	smsCode := bm.Get(req.Mobile)
	if smsCode == "" {
		beego.Info("数据获取失败")
		rsp.Errmsg = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}

	smsCodeStr, _ := redis.String(smsCode, nil)
	if req.SmsCode != smsCodeStr {
		beego.Info("短信验证码错误")
		rsp.Erron = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
		return nil
	}

	// 将数据存入数据库
	o := orm.NewOrm()
	user := models.User{
		Password_hash: Md5String(req.Password),
		Mobile:        req.Mobile,
		Name:          req.Mobile,
	}
	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("数据存储失败：", err)
		rsp.Errmsg = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return nil
	}
	beego.Info("用户ID： ", id)

	// 创建sessionid （唯一随机码）
	sessionID := Md5String(req.Mobile + req.Password)
	rsp.SessionId = sessionID

	// 以sessionid为key的一部分 创建session
	bm.Put(sessionID+"name", req.Mobile, time.Second*3600)
	bm.Put(sessionID+"userID", id, time.Second*3600)
	bm.Put(sessionID+"mobile", user.Mobile, time.Second*3600)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostRet) Stream(ctx context.Context, req *PostRegiste.StreamingRequest, stream PostRegiste.PostRet_StreamStream) error {
	log.Logf("Received PostRet.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PostRegiste.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostRet) PingPong(ctx context.Context, stream PostRegiste.PostRet_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PostRegiste.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
