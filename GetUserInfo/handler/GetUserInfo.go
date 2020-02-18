package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	GetUserInfo "sss/GetUserInfo/proto/GetUserInfo"
)

type GetUserInfo struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetUserInfo) Call(ctx context.Context, req *GetUserInfo.Request, rsp *GetUserInfo.Response) error {
	log.Log("Received GetUserInfo.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetUserInfo) Stream(ctx context.Context, req *GetUserInfo.StreamingRequest, stream GetUserInfo.GetUserInfo_StreamStream) error {
	log.Logf("Received GetUserInfo.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&GetUserInfo.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetUserInfo) PingPong(ctx context.Context, stream GetUserInfo.GetUserInfo_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&GetUserInfo.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
