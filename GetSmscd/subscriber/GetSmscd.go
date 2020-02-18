package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	GetSms "sss/GetSmscd/proto/GetSmscd"
)

type GetSmscd struct{}

func (e *GetSmscd) Handle(ctx context.Context, msg *GetSms.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *GetSms.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
