package handler

import (
	"context"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"path"
	"sss/IhomeWeb/models"
	"sss/IhomeWeb/utils"

	"github.com/micro/go-micro/util/log"

	PostAvatar_ "sss/PostAvatar/proto/PostAvatar"
)

type PostAvatar struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostAvatar) Call(ctx context.Context, req *PostAvatar_.Request, rsp *PostAvatar_.Response) error {
	beego.Info("上传用户头像 url: /api/v1.0/user/avatar")
	//初始化返回正确的返回值
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	fileext := path.Ext(req.FileName)
	beego.Info("----- ", fileext)
	resp_path, err := utils.UploadByBuffer(req.Avatar, fileext[1:])
	if err != nil {
		beego.Info("Postupavatar  models.UploadByBuffer err: ", err)
		rsp.Errno = utils.RECODE_IOERR
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

	sessionIDUserID := req.SessionID + "userID"
	valueID := bm.Get(sessionIDUserID)
	id, _ := redis.Int(valueID, nil)

	user := models.User{Id: id, Avatar_url: resp_path}
	o := orm.NewOrm()
	_, err = o.Update(&user, "avatar_url")
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
	}

	rsp.AvatarUrl = resp_path

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *PostAvatar) Stream(ctx context.Context, req *PostAvatar_.StreamingRequest, stream PostAvatar_.PostAvatar_StreamStream) error {
	log.Logf("Received PostAvatar.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&PostAvatar_.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *PostAvatar) PingPong(ctx context.Context, stream PostAvatar_.PostAvatar_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&PostAvatar_.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
