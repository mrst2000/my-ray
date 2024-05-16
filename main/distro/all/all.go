package all

import (
	// The following are necessary as they register handlers in their init functions.

	// Mandatory features. Can't remove unless there are replacements.
	_ "github.com/mrst2000/my-ray/app/dispatcher"
	_ "github.com/mrst2000/my-ray/app/proxyman/inbound"
	_ "github.com/mrst2000/my-ray/app/proxyman/outbound"

	// Default commander and all its services. This is an optional feature.
	_ "github.com/mrst2000/my-ray/app/commander"
	_ "github.com/mrst2000/my-ray/app/log/command"
	_ "github.com/mrst2000/my-ray/app/proxyman/command"
	_ "github.com/mrst2000/my-ray/app/stats/command"

	// Developer preview services
	_ "github.com/mrst2000/my-ray/app/observatory/command"

	// Other optional features.
	_ "github.com/mrst2000/my-ray/app/dns"
	_ "github.com/mrst2000/my-ray/app/dns/fakedns"
	_ "github.com/mrst2000/my-ray/app/log"
	_ "github.com/mrst2000/my-ray/app/metrics"
	_ "github.com/mrst2000/my-ray/app/policy"
	_ "github.com/mrst2000/my-ray/app/reverse"
	_ "github.com/mrst2000/my-ray/app/router"
	_ "github.com/mrst2000/my-ray/app/stats"

	// Fix dependency cycle caused by core import in internet package
	_ "github.com/mrst2000/my-ray/transport/internet/tagged/taggedimpl"

	// Developer preview features
	_ "github.com/mrst2000/my-ray/app/observatory"

	// Inbound and outbound proxies.
	_ "github.com/mrst2000/my-ray/proxy/blackhole"
	_ "github.com/mrst2000/my-ray/proxy/dns"
	_ "github.com/mrst2000/my-ray/proxy/dokodemo"
	_ "github.com/mrst2000/my-ray/proxy/freedom"
	_ "github.com/mrst2000/my-ray/proxy/http"
	_ "github.com/mrst2000/my-ray/proxy/loopback"
	_ "github.com/mrst2000/my-ray/proxy/shadowsocks"
	_ "github.com/mrst2000/my-ray/proxy/socks"
	_ "github.com/mrst2000/my-ray/proxy/trojan"
	_ "github.com/mrst2000/my-ray/proxy/vless/inbound"
	_ "github.com/mrst2000/my-ray/proxy/vless/outbound"
	_ "github.com/mrst2000/my-ray/proxy/vmess/inbound"
	_ "github.com/mrst2000/my-ray/proxy/vmess/outbound"
	_ "github.com/mrst2000/my-ray/proxy/wireguard"

	// Transports
	_ "github.com/mrst2000/my-ray/transport/internet/domainsocket"
	_ "github.com/mrst2000/my-ray/transport/internet/grpc"
	_ "github.com/mrst2000/my-ray/transport/internet/http"
	_ "github.com/mrst2000/my-ray/transport/internet/httpupgrade"
	_ "github.com/mrst2000/my-ray/transport/internet/kcp"
	_ "github.com/mrst2000/my-ray/transport/internet/quic"
	_ "github.com/mrst2000/my-ray/transport/internet/reality"
	_ "github.com/mrst2000/my-ray/transport/internet/tcp"
	_ "github.com/mrst2000/my-ray/transport/internet/tls"
	_ "github.com/mrst2000/my-ray/transport/internet/udp"
	_ "github.com/mrst2000/my-ray/transport/internet/websocket"

	// Transport headers
	_ "github.com/mrst2000/my-ray/transport/internet/headers/http"
	_ "github.com/mrst2000/my-ray/transport/internet/headers/noop"
	_ "github.com/mrst2000/my-ray/transport/internet/headers/srtp"
	_ "github.com/mrst2000/my-ray/transport/internet/headers/tls"
	_ "github.com/mrst2000/my-ray/transport/internet/headers/utp"
	_ "github.com/mrst2000/my-ray/transport/internet/headers/wechat"
	_ "github.com/mrst2000/my-ray/transport/internet/headers/wireguard"

	// JSON & TOML & YAML
	_ "github.com/mrst2000/my-ray/main/json"
	_ "github.com/mrst2000/my-ray/main/toml"
	_ "github.com/mrst2000/my-ray/main/yaml"

	// Load config from file or http(s)
	_ "github.com/mrst2000/my-ray/main/confloader/external"

	// Commands
	_ "github.com/mrst2000/my-ray/main/commands/all"
)
