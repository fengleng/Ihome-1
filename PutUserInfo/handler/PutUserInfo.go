package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"
	"time"

	"github.com/micro/go-micro/util/log"

	PutUserInfo_ "sss/PutUserInfo/proto/PutUserInfo"
)

type PutUserInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PutUserInfo) Call(ctx context.Context, req *PutUserInfo_.Request, rsp *PutUserInfo_.Response) error {
	//创建返回空间
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	value_id := bm.Get(req.SessionID + "userID")
	id, _ := redis.Int(value_id, nil)
	user := models.User{Id: id, Name: req.UserName}
	o := orm.NewOrm()
	_, err = o.Update(&user, "name")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 更新session
	sessionName := req.SessionID + "name"
	bm.Put(sessionName, user.Name, time.Second*3600)

	rsp.UserName = user.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PutUserInfo) Stream(ctx context.Context, req *PutUserInfo_.StreamingRequest, stream PutUserInfo_.PutUserInfo_StreamStream) error {
	log.Logf("Received PutUserInfo.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PutUserInfo_.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PutUserInfo) PingPong(ctx context.Context, stream PutUserInfo_.PutUserInfo_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PutUserInfo_.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
