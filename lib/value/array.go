package value

import (
	"node-go/lib/cgo"
)

type Array []Go

var _ Go = Array{}

func (v Array) Alloc() (C, error) {
	n := len(v)
	if n == 0 {
		return cgo.AllocArrayCValue(nil), nil
	}

	vs := make([]*cgo.CValue, n)
	for i := range vs {
		vi, err := v[i].Alloc()
		if err != nil {
			// try to free anything allocated so far
			for j := 0; j < i; j++ {
				vs[j].Free()
			}

			return nil, err
		}

		vs[i] = vi.(*cgo.CValue)
	}

	return cgo.AllocArrayCValue(vs), nil
}

func (Array) Kind() Kind {
	return VALUE_KIND_ARRAY
}

func (v Array) Len() int {
	return len(v)
}

func (v Array) ForEach(fn func(i int, v Go)) {
	for i, vi := range v {
		fn(i, vi)
	}
}

func (v Array) Get(i int) Go {
	if v != nil {
		return v[i]
	}

	return nil
}
