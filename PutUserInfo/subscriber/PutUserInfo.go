package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	PutUserInfo_ "sss/PutUserInfo/proto/PutUserInfo"
)

type PutUserInfo struct{}

func (e *PutUserInfo) Handle(ctx context.Context, msg *PutUserInfo_.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *PutUserInfo_.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
