package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"

	"github.com/micro/go-micro/util/log"

	GetUserInfo_ "sss/GetUserInfo/proto/GetUserInfo"
)

type GetUserInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserInfo) Call(ctx context.Context, req *GetUserInfo_.Request, rsp *GetUserInfo_.Response) error {
	beego.Info("用户信息 url: /api/v1.0/user")

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	sessionIDUserID := req.SessionID + "userID"
	userId := bm.Get(sessionIDUserID)
	id, _ := redis.Int(userId, nil)
	o := orm.NewOrm()
	user := models.User{Id: id}
	err = o.Read(&user, "id")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	rsp.UserID = int64(user.Id)
	rsp.Name = user.Name
	rsp.Mobile = user.Mobile
	rsp.RealName = user.Real_name
	rsp.IDCard = user.Id_card
	rsp.AvatarURL = user.Avatar_url

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetUserInfo) Stream(ctx context.Context, req *GetUserInfo_.StreamingRequest, stream GetUserInfo_.GetUserInfo_StreamStream) error {
	log.Logf("Received GetUserInfo.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetUserInfo_.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetUserInfo) PingPong(ctx context.Context, stream GetUserInfo_.GetUserInfo_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetUserInfo_.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
