package cgo

/*
#ifndef __LIB_FUNCTION_GO__
#define __LIB_FUNCTION_GO__

#include "../include/function.h"

#endif // __LIB_FUNCTION_GO__
*/
import "C"

import (
	"errors"
)

var (
	ErrNilFunctionContext = errors.New("function: context is nil")
)

// CFunctionContext is a wrapper around a raw C struct context.
// CFunctionContext methods will panic if the function context is invalid.
type CFunctionContext C.struct_fn_context_t

// CFunctionCallbacks is a wrapper around a raw C struct callbacks.
// CFunctionCallbacks methods do not panic if an input context is invalid.
type CFunctionCallbacks C.struct_fn_callbacks_t

func (c *CFunctionContext) IsValid() bool {
	return c != nil
}

func (c *CFunctionContext) GetCallbacks() CFunctionCallbacks {
	c.checkValid()
	return (CFunctionCallbacks)(C.function_callbacks(
		(*C.struct_fn_context_t)(c),
	))
}

func (c *CFunctionContext) NumArguments() int {
	c.checkValid()
	return int(C.function_num_arguments(
		(*C.struct_fn_context_t)(c),
	))
}

func (c *CFunctionContext) GetArgument(i int) *CValue {
	c.checkValid()
	return (*CValue)(C.function_argument(
		(*C.struct_fn_context_t)(c),
		C.size_t(i),
	))
}

func (c *CFunctionContext) Return(val *CValue) {
	c.checkValid()
	C.function_return(
		(*C.struct_fn_context_t)(c),
		(*C.struct_value_t)(val),
	)
}

func (c *CFunctionContext) Throw(err *CValue) {
	c.checkValid()
	C.function_throw(
		(*C.struct_fn_context_t)(c),
		(*C.struct_value_t)(err),
	)
}

// CFunctionContext.NewPromise allocates a new promise, if the context is valid.
func (c *CFunctionContext) NewPromise() *CPromise {
	c.checkValid()
	return (*CPromise)(C.function_new_promise(
		(*C.struct_fn_context_t)(c),
	))
}

func (cbs CFunctionCallbacks) NewPromise(ctx *CFunctionContext) *CPromise {
	// always attempt the new-promise callback, even if context is nil
	return (*CPromise)(C.function_callback_new_promise(
		(C.struct_fn_callbacks_t)(cbs),
		(*C.struct_fn_context_t)(ctx),
	))
}

// CFunctionContext.checkValid verifies that the context is valid.
// CFunctionContext.checkValid panics if the context is invalid.
func (c *CFunctionContext) checkValid() {
	if c == nil {
		panic(ErrNilFunctionContext)
	}
}
