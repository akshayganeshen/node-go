package value

import (
	"errors"
	"unsafe"

	"node-go/lib/cgo"
)

// Promise is an abstraction for a promise fulfillable.
type Promise interface {
	IsPending() bool
	Resolve(val Go) error
	Reject(err Go) error
}

// GoPromise is a lightweight wrapper around a raw cgo.CPromise.
// GoPromise is freed when it is fulfilled.
type GoPromise struct {
	Handle  *cgo.CPromise
	Pending bool
}

var _ Promise = &GoPromise{}
var _ Go = &GoPromise{}
var _ C = &GoPromise{}

var (
	ErrPromiseFulfilled = errors.New("promise: promise already fulfilled")
)

func (p *GoPromise) IsPending() bool {
	return p.Pending
}

func (p *GoPromise) Resolve(val Go) error {
	if !p.Pending {
		return ErrPromiseFulfilled
	}

	vc, err := val.Alloc()
	if err != nil {
		return err
	}

	defer p.Free()
	defer p.fulfill()
	defer vc.Free()
	return p.Handle.Resolve(vc.(*cgo.CValue))
}

func (p *GoPromise) Reject(erv Go) error {
	if !p.Pending {
		return ErrPromiseFulfilled
	}

	ec, err := erv.Alloc()
	if err != nil {
		return err
	}

	defer p.Free()
	defer p.fulfill()
	defer ec.Free()
	return p.Handle.Reject(ec.(*cgo.CValue))
}

// GoPromise.Alloc does not actually allocate a C promise struct.
// It just allocates the container value to point to it.
func (p *GoPromise) Alloc() (C, error) {
	return cgo.AllocPromiseCValue(p.Handle), nil
}

func (GoPromise) Kind() Kind {
	return VALUE_KIND_PROMISE
}

func (p *GoPromise) Free() error {
	handle := p.Handle
	p.Handle = nil
	return handle.Free()
}

func (p *GoPromise) Pointer() uintptr {
	return uintptr(unsafe.Pointer(p.Handle))
}

func (p *GoPromise) fulfill() {
	p.Pending = false
}
