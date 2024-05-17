package websocket

import (
	"net/http"
	"fmt"
	"github.com/mrst2000/my-ray/common"
	"github.com/mrst2000/my-ray/transport/internet"
)

const protocolName = "websocket"

func main() {
	config := &Config{
		Header: map[string]string{
			"Content-Type": "application/json",
			"Authorization": "Bearer token",
		},
		Host: "example.com",
	}

	header := config.GetRequestHeader()
	// Modify the header map
	header["hoSt"] = header["Host"]
	delete(header, "Host")
}

type Config struct {
	Header map[string]string
	Host   string
}

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
	randomizedHost := randomizeCase(c.Host)
	header.Set("Host", randomizedHost)
	header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36")
	return header
}

func randomizeCase(s string) string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune(s)
	for i, r := range runes {
		if rand.Intn(2) == 0 {
			runes[i] = unicode.ToLower(r)
		} else {
			runes[i] = unicode.ToUpper(r)
		}
	}
	return string(runes)
}

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
