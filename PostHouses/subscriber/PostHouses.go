package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	PostHouses "sss/PostHouses/proto/PostHouses"
)

type PostHouse struct{}

func (e *PostHouse) Handle(ctx context.Context, msg *PostHouses.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *PostHouses.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
