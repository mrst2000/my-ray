package udp

import (
	"github.com/mrst2000/my-ray/common"
	"github.com/mrst2000/my-ray/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
