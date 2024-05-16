package shadowsocks

import (
	"context"
	"time"

	"github.com/mrst2000/my-ray/common"
	"github.com/mrst2000/my-ray/common/buf"
	"github.com/mrst2000/my-ray/common/net"
	"github.com/mrst2000/my-ray/common/protocol"
	"github.com/mrst2000/my-ray/common/retry"
	"github.com/mrst2000/my-ray/common/session"
	"github.com/mrst2000/my-ray/common/signal"
	"github.com/mrst2000/my-ray/common/task"
	"github.com/mrst2000/my-ray/core"
	"github.com/mrst2000/my-ray/features/policy"
	"github.com/mrst2000/my-ray/transport"
	"github.com/mrst2000/my-ray/transport/internet"
	"github.com/mrst2000/my-ray/transport/internet/stat"
)

// Client is a inbound handler for Shadowsocks protocol
type Client struct {
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
}

// NewClient create a new Shadowsocks client.
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Server {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, newError("failed to parse server spec").Base(err)
		}
		serverList.AddServer(s)
	}
	if serverList.Size() == 0 {
		return nil, newError("0 server")
	}

	v := core.MustFromContext(ctx)
	client := &Client{
		serverPicker:  protocol.NewRoundRobinServerPicker(serverList),
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
	}
	return client, nil
}

// Process implements OutboundHandler.Process().
func (c *Client) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	outbounds := session.OutboundsFromContext(ctx)
	ob := outbounds[len(outbounds) - 1]
	if !ob.Target.IsValid() {
		return newError("target not specified")
	}
	ob.Name = "shadowsocks"
	ob.CanSpliceCopy = 3
	destination := ob.Target
	network := destination.Network

	var server *protocol.ServerSpec
	var conn stat.Connection

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		server = c.serverPicker.PickServer()
		dest := server.Destination()
		dest.Network = network
		rawConn, err := dialer.Dial(ctx, dest)
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		return newError("failed to find an available destination").AtWarning().Base(err)
	}
	newError("tunneling request to ", destination, " via ", network, ":", server.Destination().NetAddr()).WriteToLog(session.ExportIDToError(ctx))

	defer conn.Close()

	request := &protocol.RequestHeader{
		Version: Version,
		Address: destination.Address,
		Port:    destination.Port,
	}
	if destination.Network == net.Network_TCP {
		request.Command = protocol.RequestCommandTCP
	} else {
		request.Command = protocol.RequestCommandUDP
	}

	user := server.PickUser()
	_, ok := user.Account.(*MemoryAccount)
	if !ok {
		return newError("user account is not valid")
	}
	request.User = user

	var newCtx context.Context
	var newCancel context.CancelFunc
	if session.TimeoutOnlyFromContext(ctx) {
		newCtx, newCancel = context.WithCancel(context.Background())
	}

	sessionPolicy := c.policyManager.ForLevel(user.Level)
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, func() {
		cancel()
		if newCancel != nil {
			newCancel()
		}
	}, sessionPolicy.Timeouts.ConnectionIdle)

	if newCtx != nil {
		ctx = newCtx
	}

	if request.Command == protocol.RequestCommandTCP {
		requestDone := func() error {
			defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)
			bufferedWriter := buf.NewBufferedWriter(buf.NewWriter(conn))
			bodyWriter, err := WriteTCPRequest(request, bufferedWriter)
			if err != nil {
				return newError("failed to write request").Base(err)
			}

			if err = buf.CopyOnceTimeout(link.Reader, bodyWriter, time.Millisecond*100); err != nil && err != buf.ErrNotTimeoutReader && err != buf.ErrReadTimeout {
				return newError("failed to write A request payload").Base(err).AtWarning()
			}

			if err := bufferedWriter.SetBuffered(false); err != nil {
				return err
			}

			return buf.Copy(link.Reader, bodyWriter, buf.UpdateActivity(timer))
		}

		responseDone := func() error {
			defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

			responseReader, err := ReadTCPResponse(user, conn)
			if err != nil {
				return err
			}

			return buf.Copy(responseReader, link.Writer, buf.UpdateActivity(timer))
		}

		responseDoneAndCloseWriter := task.OnSuccess(responseDone, task.Close(link.Writer))
		if err := task.Run(ctx, requestDone, responseDoneAndCloseWriter); err != nil {
			return newError("connection ends").Base(err)
		}

		return nil
	}

	if request.Command == protocol.RequestCommandUDP {

		requestDone := func() error {
			defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

			writer := &UDPWriter{
				Writer:  conn,
				Request: request,
			}

			if err := buf.Copy(link.Reader, writer, buf.UpdateActivity(timer)); err != nil {
				return newError("failed to transport all UDP request").Base(err)
			}
			return nil
		}

		responseDone := func() error {
			defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

			reader := &UDPReader{
				Reader: conn,
				User:   user,
			}

			if err := buf.Copy(reader, link.Writer, buf.UpdateActivity(timer)); err != nil {
				return newError("failed to transport all UDP response").Base(err)
			}
			return nil
		}

		responseDoneAndCloseWriter := task.OnSuccess(responseDone, task.Close(link.Writer))
		if err := task.Run(ctx, requestDone, responseDoneAndCloseWriter); err != nil {
			return newError("connection ends").Base(err)
		}

		return nil
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
