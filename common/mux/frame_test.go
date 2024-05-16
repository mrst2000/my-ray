package mux_test

import (
	"testing"

	"github.com/mrst2000/my-ray/common"
	"github.com/mrst2000/my-ray/common/buf"
	"github.com/mrst2000/my-ray/common/mux"
	"github.com/mrst2000/my-ray/common/net"
)

func BenchmarkFrameWrite(b *testing.B) {
	frame := mux.FrameMetadata{
		Target:        net.TCPDestination(net.DomainAddress("www.example.com"), net.Port(80)),
		SessionID:     1,
		SessionStatus: mux.SessionStatusNew,
	}
	writer := buf.New()
	defer writer.Release()

	for i := 0; i < b.N; i++ {
		common.Must(frame.WriteTo(writer))
		writer.Clear()
	}
}
