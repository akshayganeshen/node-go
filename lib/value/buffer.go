package value

import (
	"node-go/lib/cgo"
)

type Buffer []byte

var _ Go = Buffer{}

func (b Buffer) Alloc() (C, error) {
	return cgo.AllocBufferCValue([]byte(b)), nil
}

func (Buffer) Kind() Kind {
	return VALUE_KIND_BUFFER
}
