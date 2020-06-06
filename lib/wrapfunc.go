package main

/*
#ifndef __LIB_WRAPFUNC_GO__
#define __LIB_WRAPFUNC_GO__

#include "./include/function.h"

#endif // __LIB_WRAPFUNC_GO__
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	"node-go/lib/cgo"
	"node-go/lib/function"
)

type CGoHandler func(req *C.struct_fn_context_t)
type NodeGoHandler func(ctx function.Context) error

func WrapHandler(h NodeGoHandler) CGoHandler {
	return func(req *C.struct_fn_context_t) {
		ctx := BuildContext(req)
		defer ctx.ClearContextHandle()

		if err := h(ctx); err != nil {
			// TODO: Automatically convert error and to use throw callback.
			fmt.Fprintf(os.Stderr, "wrapfunc: handler error: %s\n", err)
		}
	}
}

func BuildContext(req *C.struct_fn_context_t) *function.GoContext {
	return function.CopyC(
		cgo.UnsafeFrom(unsafe.Pointer(req)).ToCFunctionContext(),
	)
}
