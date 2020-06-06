package value

import (
	"node-go/lib/cgo"
)

type Number float64

var _ Go = Number(0)

func (v Number) Alloc() (C, error) {
	return cgo.AllocNumberCValue(float64(v)), nil
}

func (Number) Kind() Kind {
	return VALUE_KIND_NUMBER
}
