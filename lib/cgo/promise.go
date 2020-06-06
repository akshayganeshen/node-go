package cgo

/*
#ifndef __LIB_PROMISE_GO__
#define __LIB_PROMISE_GO__

#include "../include/promise.h"

#endif // __LIB_PROMISE_GO__
*/
import "C"

// CPromise is a wrapper around a raw C struct promise.
// CPromise pointers are expected to be managed by Go.
type CPromise C.struct_promise_t

func (p *CPromise) Resolve(val *CValue) error {
	C.promise_resolve(
		(*C.struct_promise_t)(p),
		(*C.struct_value_t)(val),
	)
	return nil
}

func (p *CPromise) Reject(err *CValue) error {
	C.promise_reject(
		(*C.struct_promise_t)(p),
		(*C.struct_value_t)(err),
	)
	return nil
}

func (p *CPromise) Free() error {
	C.promise_free(
		(*C.struct_promise_t)(p),
	)
	return nil
}
