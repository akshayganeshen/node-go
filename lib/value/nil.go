package value

import (
	"node-go/lib/cgo"
)

type Null struct{}
type Undefined struct{}

var _ Go = Null{}
var _ Go = Undefined{}

func (Null) Alloc() (C, error) {
	return cgo.AllocNullCValue(), nil
}

func (Null) Kind() Kind {
	return VALUE_KIND_NULL
}

func (Undefined) Alloc() (C, error) {
	return cgo.AllocUndefinedCValue(), nil
}

func (Undefined) Kind() Kind {
	return VALUE_KIND_UNDEFINED
}
