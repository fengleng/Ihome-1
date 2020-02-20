package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	GetHouses_ "sss/GetHouses/proto/GetHouses"
)

type GetHouses struct{}

func (e *GetHouses) Handle(ctx context.Context, msg *GetHouses_.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *GetHouses_.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
