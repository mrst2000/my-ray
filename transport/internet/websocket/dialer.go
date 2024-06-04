package websocket

import (
    "bytes"
    "context"
    _ "embed"
    "encoding/base64"
    "io"
    gonet "net"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
    "github.com/xtls/xray-core/common"
    "github.com/xtls/xray-core/common/net"
    "github.com/xtls/xray-core/common/platform"
    "github.com/xtls/xray-core/common/session"
    "github.com/xtls/xray-core/common/uuid"
    "github.com/xtls/xray-core/transport/internet"
    "github.com/xtls/xray-core/transport/internet/stat"
    "github.com/xtls/xray-core/transport/internet/tls"
)

//go:embed dialer.html
var webpage []byte

var conns chan *websocket.Conn

func init() {
    addr := platform.NewEnvFlag(platform.BrowserDialerAddress).GetValue(func() string { return "" })
    if addr != "" {
        token := uuid.New()
        csrfToken := token.String()
        webpage = bytes.ReplaceAll(webpage, []byte("csrfToken"), []byte(csrfToken))
        conns = make(chan *websocket.Conn, 256)
        go http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.URL.Path == "/websocket" {
                if r.URL.Query().Get("token") == csrfToken {
                    if conn, err := upgrader.Upgrade(w, r, nil); err == nil {
                        conns <- conn
                    } else {
                        newError("Browser dialer http upgrade unexpected error").AtError().WriteToLog()
                    }
                }
            } else {
                w.Write(webpage)
            }
        }))
    }
}

// Dial dials a WebSocket connection to the given destination.
func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (stat.Connection, error) {
    newError("creating connection to ", dest).WriteToLog(session.ExportIDToError(ctx))
    var conn net.Conn
    if streamSettings.ProtocolSettings.(*Config).Ed > 0 {
        ctx, cancel := context.WithCancel(ctx)
        conn = &delayDialConn{
            dialed:         make(chan bool, 1),
            cancel:         cancel,
            ctx:            ctx,
            dest:           dest,
            streamSettings: streamSettings,
        }
    } else {
        var err error
        conn, err = internet.DialSystem(ctx, dest, streamSettings)
        if err != nil {
            return nil, err
        }

        // Send the fake TLS client hello
        err = tls.SendFakeTLSClientHello(conn, "www.speedtest.net")
        if err != nil {
            conn.Close()
            return nil, err
        }

        // Now proceed with the actual TLS handshake
        tlsConfig := &tls.Config{
            ServerName: streamSettings.ProtocolSettings.(*Config).Host,
        }
        tlsConn := tls.Client(conn, tlsConfig)
        err = tlsConn.Handshake()
        if err != nil {
            conn.Close()
            return nil, err
        }
        conn = tlsConn
    }

    header := streamSettings.ProtocolSettings.(*Config).GetRequestHeader()
    uri := "wss://" + dest.Address.String() + streamSettings.ProtocolSettings.(*Config).GetNormalizedPath()

    dialer := websocket.Dialer{
        Proxy:            http.ProxyFromEnvironment,
        HandshakeTimeout: 45 * time.Second,
        TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
    }

    if streamSettings.ProtocolSettings.(*Config).Ed != nil {
        header.Set("Sec-WebSocket-Protocol", base64.RawURLEncoding.EncodeToString(streamSettings.ProtocolSettings.(*Config).Ed))
    }

    wsConn, resp, err := dialer.DialContext(ctx, uri, header)
    if err != nil {
        var reason string
        if resp != nil {
            reason = resp.Status
        }
        return nil, newError("failed to dial to (", uri, "): ", reason).Base(err)
    }

    return newConnection(wsConn, wsConn.RemoteAddr(), nil), nil
}

type delayDialConn struct {
    net.Conn
    closed         bool
    dialed         chan bool
    cancel         context.CancelFunc
    ctx            context.Context
    dest           net.Destination
    streamSettings *internet.MemoryStreamConfig
}

func (d *delayDialConn) Write(b []byte) (int, error) {
    if d.closed {
        return 0, io.ErrClosedPipe
    }
    if d.Conn == nil {
        ed := b
        if len(ed) > int(d.streamSettings.ProtocolSettings.(*Config).Ed) {
            ed = nil
        }
        var err error
        if d.Conn, err = dialWebSocket(d.ctx, d.dest, d.streamSettings, ed); err != nil {
            d.Close()
            return 0, newError("failed to dial WebSocket").Base(err)
        }
        d.dialed <- true
        if ed != nil {
            return len(ed), nil
        }
    }
    return d.Conn.Write(b)
}

func (d *delayDialConn) Read(b []byte) (int, error) {
    if d.closed {
        return 0, io.ErrClosedPipe
    }
    if d.Conn == nil {
        select {
        case <-d.ctx.Done():
            return 0, io.ErrUnexpectedEOF
        case <-d.dialed:
        }
    }
    return d.Conn.Read(b)
}

func (d *delayDialConn) Close() error {
    d.closed = true
    d.cancel()
    if d.Conn == nil {
        return nil
    }
    return d.Conn.Close()
}
