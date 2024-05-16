package drain

import "io"

//go:generate go run github.com/mrst2000/my-ray/common/errors/errorgen

type Drainer interface {
	AcknowledgeReceive(size int)
	Drain(reader io.Reader) error
}
