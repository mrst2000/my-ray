package dokodemo

import (
	"github.com/mrst2000/my-ray/common/net"
)

// GetPredefinedAddress returns the defined address from proto config. Null if address is not valid.
func (v *Config) GetPredefinedAddress() net.Address {
	addr := v.Address.AsAddress()
	if addr == nil {
		return nil
	}
	return addr
}
