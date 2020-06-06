package value

import (
	"fmt"
)

type Object map[string]Go

var _ Go = Object{}

func (v Object) Alloc() (C, error) {
	return nil, fmt.Errorf("Object.Alloc: unimplemented")
}

func (Object) Kind() Kind {
	return VALUE_KIND_OBJECT
}

func (v Object) Len() int {
	return len(v)
}

func (v Object) ForEach(fn func(key string, val Go)) {
	if v != nil {
		for key, value := range v {
			fn(key, value)
		}
	}
}
