package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	GetIndex_ "sss/GetIndex/proto/GetIndex"
)

type GetIndex struct{}

func (e *GetIndex) Handle(ctx context.Context, msg *GetIndex_.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *GetIndex_.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
