package conf

import (
	"github.com/mrst2000/my-ray/proxy/loopback"
	"google.golang.org/protobuf/proto"
)

type LoopbackConfig struct {
	InboundTag string `json:"inboundTag"`
}

func (l LoopbackConfig) Build() (proto.Message, error) {
	return &loopback.Config{InboundTag: l.InboundTag}, nil
}
