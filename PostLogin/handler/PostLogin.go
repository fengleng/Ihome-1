package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"sss/PostRet/handler"
	"time"

	"github.com/micro/go-micro/util/log"

	PostLogin_ "sss/PostLogin/proto/PostLogin"
)

type PostLogin struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostLogin) Call(ctx context.Context, req *PostLogin_.Request, rsp *PostLogin_.Response) error {
	beego.Info("用户登陆 url: /api/v1.0/sessions")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	// 用户验证是否存在
	user := models.User{Mobile: req.Mobile}
	o := orm.NewOrm()
	err := o.Read(&user, "mobile")
	beego.Info("数据库查找：", err)
	if err != nil {
		beego.Info("用户不存在")
		rsp.Errno = utils.RECODE_USERERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	// 密码验证
	beego.Info("数据库密码：", user.Password_hash, ", 输入密码：", handler.Md5String(req.Password), req.Mobile)
	if user.Password_hash != handler.Md5String(req.Password) {
		beego.Info("用户名或者密码不正确")
		rsp.Errno = utils.RECODE_PWDERR
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

	// 创建sessionid （唯一随机码）
	sessionID := handler.Md5String(req.Mobile + req.Password)
	rsp.SessionID = sessionID

	// 以sessionid为key的一部分 创建session
	bm.Put(sessionID+"name", req.Mobile, time.Second*3600)
	bm.Put(sessionID+"userID", user.Id, time.Second*3600)
	bm.Put(sessionID+"mobile", user.Mobile, time.Second*3600)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostLogin) Stream(ctx context.Context, req *PostLogin_.StreamingRequest, stream PostLogin_.PostLogin_StreamStream) error {
	log.Logf("Received PostLogin.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PostLogin_.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostLogin) PingPong(ctx context.Context, stream PostLogin_.PostLogin_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PostLogin_.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
