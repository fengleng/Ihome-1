package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	GetUserInfo_ "sss/GetUserInfo/proto/GetUserInfo"
)

type GetUserInfo struct{}

func (e *GetUserInfo) Handle(ctx context.Context, msg *GetUserInfo_.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *GetUserInfo_.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
