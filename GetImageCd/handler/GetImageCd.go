package handler

import (
	"context"
	"github.com/afocus/captcha"
	"github.com/astaxie/beego"
	"image/color"
	"sss/IhomeWeb/utils"
	"time"

	"github.com/micro/go-micro/util/log"
	GetImage "sss/GetImageCd/proto/GetImageCd"

	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

type GetImageCd struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetImageCd) Call(ctx context.Context, req *GetImage.Request, rsp *GetImage.Response) error {
	beego.Info("获取首页轮播图 url: /api/v1.0/imagecode")

	cap := generateImage()

	img, str := cap.Create(4, captcha.ALL)

	bm, err := utils.Redis(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		beego.Info("redis 连接失败: ", err)
		rsp.Error = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Error)
	}
	// 将验证码和uuid存入redis中
	bm.Put(req.Uuid, str, time.Second*300)

	rsp.Error = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

	img1 := *img
	img_t := *img1.RGBA
	rsp.Pix = img_t.Pix
	rsp.Stride = int64(img_t.Stride)
	rsp.Min = &GetImage.Response_Point{
		X: int64(img_t.Rect.Min.X),
		Y: int64(img_t.Rect.Min.Y),
	}
	rsp.Max = &GetImage.Response_Point{
		X: int64(img_t.Rect.Max.X),
		Y: int64(img_t.Rect.Max.Y),
	}

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetImageCd) Stream(ctx context.Context, req *GetImage.StreamingRequest, stream GetImage.GetImageCd_StreamStream) error {
	log.Logf("Received GetImageCd.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetImage.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetImageCd) PingPong(ctx context.Context, stream GetImage.GetImageCd_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetImage.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

// 生成验证码图片
func generateImage() *captcha.Captcha {
	cap := captcha.New()

	if err := cap.SetFont("comic.ttf"); err != nil {
		panic(err.Error())
	}

	cap.SetSize(90, 41)
	cap.SetDisturbance(captcha.MEDIUM)
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	return cap
}
