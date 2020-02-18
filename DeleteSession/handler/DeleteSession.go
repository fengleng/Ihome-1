package handler

import (
	"context"
	"github.com/astaxie/beego"
	"sss/IhomeWeb/utils"

	"github.com/micro/go-micro/util/log"

	DeleteSession_ "sss/DeleteSession/proto/DeleteSession"
)

type DeleteSession struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *DeleteSession) Call(ctx context.Context, req *DeleteSession_.Request, rsp *DeleteSession_.Response) error {
	beego.Info("用户登陆 url: /api/v1.0/session")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)
	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	sessionIdName := req.SessionID + "name"
	sessionIdUserID := req.SessionID + "userID"
	sessionIdMobile := req.SessionID + "mobile"
	bm.Delete(sessionIdName)
	bm.Delete(sessionIdMobile)
	bm.Delete(sessionIdUserID)

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *DeleteSession) Stream(ctx context.Context, req *DeleteSession_.StreamingRequest, stream DeleteSession_.DeleteSession_StreamStream) error {
	log.Logf("Received DeleteSession.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&DeleteSession_.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *DeleteSession) PingPong(ctx context.Context, stream DeleteSession_.DeleteSession_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&DeleteSession_.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
