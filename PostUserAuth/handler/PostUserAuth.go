package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"

	"github.com/micro/go-micro/util/log"

	PostUserAuth_ "sss/PostUserAuth/proto/PostUserAuth"
)

type PostUserAuth struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostUserAuth) Call(ctx context.Context, req *PostUserAuth_.Request, rsp *PostUserAuth_.Response) error {
	beego.Info("更新实名认证 url: /api/v1.0/user/auth")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

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
	user := models.User{Id: id, Real_name: req.RealName, Id_card: req.IDCard}
	o := orm.NewOrm()
	_, err = o.Update(&user, "real_name", "id_card")
	if err != nil {
		beego.Info("用户信息更新失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostUserAuth) Stream(ctx context.Context, req *PostUserAuth_.StreamingRequest, stream PostUserAuth_.PostUserAuth_StreamStream) error {
	log.Logf("Received PostUserAuth.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PostUserAuth_.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostUserAuth) PingPong(ctx context.Context, stream PostUserAuth_.PostUserAuth_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PostUserAuth_.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
