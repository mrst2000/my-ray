package websocket

import (
	"net/http"

	"github.com/xtls/xray-core/common"
	utls "github.com/refraction-networking/utls"
	"github.com/xtls/xray-core/transport/internet"
)

const protocolName = "websocket"

func (c *Config) GetNormalizedPath() string {
	path := c.Path
	if path == "" {
		return "/"
	}
	if path[0] != '/' {
		return "/" + path
	}
	return path
}

func (c *Config) GetRequestHeader() http.Header {
	header := http.Header{}
	for k, v := range c.Header {
		header.Add(k, v)
	}
	header.Set("hoSt", c.Host)
	header.Set("User-Agent", string(utls.HelloChrome_Auto))
	return header
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
