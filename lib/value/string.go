package value

import (
	"node-go/lib/cgo"
)

type String string

var _ Go = String("")

func (s String) Alloc() (C, error) {
	return cgo.AllocStringCValue(string(s)), nil
}

func (String) Kind() Kind {
	return VALUE_KIND_STRING
}
