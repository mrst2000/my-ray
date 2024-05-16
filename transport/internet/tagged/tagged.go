package tagged

import (
	"context"

	"github.com/mrst2000/my-ray/common/net"
)

type DialFunc func(ctx context.Context, dest net.Destination, tag string) (net.Conn, error)

var Dialer DialFunc
