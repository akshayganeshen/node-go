package value

import (
	"node-go/lib/cgo"
)

type Boolean bool

var _ Go = Boolean(false)

func (v Boolean) Alloc() (C, error) {
	return cgo.AllocBooleanCValue(bool(v)), nil
}

func (Boolean) Kind() Kind {
	return VALUE_KIND_BOOLEAN
}
