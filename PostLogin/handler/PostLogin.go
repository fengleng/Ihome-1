package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	PostLogin_ "sss/PostLogin/proto/PostLogin"
)

type PostLogin struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *PostLogin) Call(ctx context.Context, req *PostLogin_.Request, rsp *PostLogin_.Response) error {
	log.Log("Received PostLogin.Call request")
	rsp.Msg = "Hello " + req.Name
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
