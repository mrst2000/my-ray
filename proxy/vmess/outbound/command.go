package outbound

import (
	"time"

	"github.com/mrst2000/my-ray/common"
	"github.com/mrst2000/my-ray/common/net"
	"github.com/mrst2000/my-ray/common/protocol"
	"github.com/mrst2000/my-ray/proxy/vmess"
)

func (h *Handler) handleSwitchAccount(cmd *protocol.CommandSwitchAccount) {
	rawAccount := &vmess.Account{
		Id: cmd.ID.String(),
		SecuritySettings: &protocol.SecurityConfig{
			Type: protocol.SecurityType_AUTO,
		},
	}

	account, err := rawAccount.AsAccount()
	common.Must(err)
	user := &protocol.MemoryUser{
		Email:   "",
		Level:   cmd.Level,
		Account: account,
	}
	dest := net.TCPDestination(cmd.Host, cmd.Port)
	until := time.Now().Add(time.Duration(cmd.ValidMin) * time.Minute)
	h.serverList.AddServer(protocol.NewServerSpec(dest, protocol.BeforeTime(until), user))
}

func (h *Handler) handleCommand(dest net.Destination, cmd protocol.ResponseCommand) {
	switch typedCommand := cmd.(type) {
	case *protocol.CommandSwitchAccount:
		if typedCommand.Host == nil {
			typedCommand.Host = dest.Address
		}
		h.handleSwitchAccount(typedCommand)
	default:
	}
}
