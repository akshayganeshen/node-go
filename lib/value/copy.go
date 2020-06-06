package value

import (
	"fmt"
	"os"

	"node-go/lib/cgo"
)

// CopyC creates a Go copy of a cgo.CValue.
// CopyC returns a Go pointer, so it is garbage collected automatically.
func CopyC(val *cgo.CValue) Go {
	if val.IsNull() {
		return Null{}
	}

	if val.IsUndefined() {
		return Undefined{}
	}

	if val.IsBoolean() {
		return Boolean(val.GetBoolean())
	}

	if val.IsNumber() {
		return Number(val.GetNumber())
	}

	if val.IsString() {
		return String(val.GetString())
	}

	if val.IsBuffer() {
		return Buffer(val.GetBytes())
	}

	if val.IsArray() {
		n := val.ArrayLen()
		gv := make([]Go, n)

		val.ArrayForEach(func(i int, v *cgo.CValue) {
			gv[i] = CopyC(v)
		})

		return Array(gv)
	}

	if val.IsPromise() {
		return CopyCPromise(val.GetPromise())
	}

	fmt.Fprintf(os.Stderr, "value: copy: unsupported value: %+v\n", val)
	return nil
}

func CopyCPromise(prom *cgo.CPromise) *GoPromise {
	return &GoPromise{
		Handle:  prom,
		Pending: true,
	}
}
