package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"github.com/micro/go-micro/util/log"
	"sss/IhomeWeb/utils"

	GetSessions "sss/GetSession/proto/GetSession"
)

type GetSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetSession) Call(ctx context.Context, req *GetSessions.Request, rsp *GetSessions.Response) error {
	beego.Info("获取Session url:/api/v1.0/session")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	// 获取username
	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	name := bm.Get(req.SessionID + "name")
	// 没有返回失败
	if name == nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	// 有则返回成功
	username, _ := redis.String(name, nil)
	rsp.UserName = username

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetSession) Stream(ctx context.Context, req *GetSessions.StreamingRequest, stream GetSessions.GetSession_StreamStream) error {
	log.Logf("Received GetSession.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetSessions.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetSession) PingPong(ctx context.Context, stream GetSessions.GetSession_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetSessions.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
