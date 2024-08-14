package biz

import "github.com/ssonit/aura_server/internal/pin/utils"

type service struct {
	store utils.PinStore
}

func NewService(store utils.PinStore) *service {
	return &service{store: store}
}
