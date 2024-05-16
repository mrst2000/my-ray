package all

import (
	"github.com/mrst2000/my-ray/main/commands/all/api"
	"github.com/mrst2000/my-ray/main/commands/all/tls"
	"github.com/mrst2000/my-ray/main/commands/base"
)

// go:generate go run github.com/mrst2000/my-ray/common/errors/errorgen

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		// cmdConvert,
		tls.CmdTLS,
		cmdUUID,
		cmdX25519,
		cmdWG,
	)
}
