package function

import (
	"errors"

	"node-go/lib/cgo"
	"node-go/lib/value"
)

var (
	ErrFunctionReturned = errors.New("function: function has already returned")
)

type Context interface {
	NumArguments() int
	GetArgument(i int) value.Js
	Callbacks() Callbacks
}

type Callbacks interface {
	Return(val value.Go) error
	Throw(err value.Go) error
	NewPromise() (value.Promise, error)
	// no DeletePromise, since promises are able to free themselves
}

type GoContext struct {
	ContextHandle *cgo.CFunctionContext
	CallbacksCopy cgo.CFunctionCallbacks
	Arguments     []value.Js
	Running       bool
}

var _ Context = &GoContext{}
var _ Callbacks = &GoContext{}

func (c GoContext) NumArguments() int {
	return len(c.Arguments)
}

func (c GoContext) GetArgument(i int) value.Js {
	return c.Arguments[i]
}

func (c *GoContext) Callbacks() Callbacks {
	return c
}

func (c *GoContext) Return(val value.Go) error {
	if !c.Running {
		return ErrFunctionReturned
	}

	vc, err := val.Alloc()
	if err != nil {
		return err
	}

	defer c.finish()
	c.ContextHandle.Return(vc.(*cgo.CValue))
	return nil
}

func (c *GoContext) Throw(erv value.Go) error {
	if !c.Running {
		return ErrFunctionReturned
	}

	ec, err := erv.Alloc()
	if err != nil {
		return err
	}

	defer c.finish()
	c.ContextHandle.Throw(ec.(*cgo.CValue))
	return nil
}

func (c *GoContext) NewPromise() (value.Promise, error) {
	// attempt the allocation even if the function has returned
	// the callbacks copy ensures we still have a reference to it
	cp := c.CallbacksCopy.NewPromise(c.ContextHandle)

	if cp != nil {
		return value.CopyCPromise(cp), nil
	}

	if !c.ContextHandle.IsValid() {
		// remap nil context to error for convenience
		return nil, ErrFunctionReturned
	}

	return nil, nil
}

func (c *GoContext) ClearContextHandle() error {
	defer c.finish()
	c.ContextHandle = nil
	return nil
}

func (c *GoContext) finish() {
	c.Running = false
}
