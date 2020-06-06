package cgo

/*
#ifndef __LIB_VALUE_GO__
#define __LIB_VALUE_GO__

#include "../include/value.h"

static inline struct value_t *bridge_alloc_value_string(
  _GoString_ str
) {
  return alloc_value_string(
    _GoStringPtr(str),
    _GoStringLen(str)
  );
}

static inline struct value_t *bridge_get_array_item(
  struct value_t **items,
  size_t i
) {
  return items[i];
}

#endif // __LIB_VALUE_GO__
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// CValue is a wrapper around a raw C struct value.
// CValue pointers are expected to be managed by C.
// CValue implements the value.C interface.
type CValue C.struct_value_t

type UnknownValueKindError struct {
	Kind int
}

func (e UnknownValueKindError) Error() string {
	return fmt.Sprintf("value: unknown kind '%d'", e.Kind)
}

func AllocUndefinedCValue() *CValue {
	return (*CValue)(C.alloc_value_undefined())
}

func AllocNullCValue() *CValue {
	return (*CValue)(C.alloc_value_null())
}

func AllocBooleanCValue(bol bool) *CValue {
	cb := C.uchar(0)
	if bol {
		cb = C.uchar(1)
	}
	return (*CValue)(C.alloc_value_boolean(cb))
}

func AllocNumberCValue(num float64) *CValue {
	return (*CValue)(C.alloc_value_number(C.double(num)))
}

func AllocStringCValue(str string) *CValue {
	return (*CValue)(C.bridge_alloc_value_string(str))
}

func AllocBufferCValue(buf []byte) *CValue {
	return (*CValue)(C.alloc_value_buffer(
		unsafe.Pointer(&buf[0]),
		C.size_t(len(buf)),
	))
}

func AllocArrayCValue(items []*CValue) *CValue {
	n := len(items)
	if n == 0 {
		return (*CValue)(C.alloc_value_array(
			nil,
			C.size_t(len(items)),
		))
	}

	return (*CValue)(C.alloc_value_array(
		(**C.struct_value_t)(unsafe.Pointer(&items[0])),
		C.size_t(n),
	))
}

func AllocPromiseCValue(prom *CPromise) *CValue {
	return (*CValue)(C.alloc_value_promise(
		(*C.struct_promise_t)(prom),
	))
}

func (v CValue) IsNull() bool {
	return v.kind == C.VALUE_KIND_NULL
}

func (v CValue) IsUndefined() bool {
	return v.kind == C.VALUE_KIND_UNDEFINED
}

func (v CValue) IsNil() bool {
	return v.IsNull() || v.IsUndefined()
}

func (v CValue) IsError() bool {
	return v.kind == C.VALUE_KIND_ERROR
}

func (v CValue) IsBoolean() bool {
	return v.kind == C.VALUE_KIND_BOOLEAN
}

func (v CValue) IsNumber() bool {
	return v.kind == C.VALUE_KIND_NUMBER
}

func (v CValue) IsString() bool {
	return v.kind == C.VALUE_KIND_STRING
}

func (v CValue) IsBuffer() bool {
	return v.kind == C.VALUE_KIND_BUFFER
}

func (v CValue) IsStringOrBuffer() bool {
	return v.IsString() || v.IsBuffer()
}

func (v CValue) IsArray() bool {
	return v.kind == C.VALUE_KIND_ARRAY
}

func (v CValue) IsObject() bool {
	return v.kind == C.VALUE_KIND_OBJECT
}

func (v CValue) IsPromise() bool {
	return v.kind == C.VALUE_KIND_PROMISE
}

func (v CValue) GetBoolean() bool {
	// NOTE: Not valid to call this on non-boolean values.
	if v.kind != C.VALUE_KIND_BOOLEAN {
		return false
	}

	return v.boolean != 0
}

func (v CValue) GetNumber() float64 {
	if v.kind != C.VALUE_KIND_NUMBER {
		return 0
	}

	return float64(v.number)
}

func (v CValue) GetString() string {
	switch v.kind {
	default:
		return ""

	case C.VALUE_KIND_ERROR:
		if v.error.reason != nil {
			return C.GoStringN(
				v.error.reason.data, C.int(v.error.reason.len),
			)
		}

		return ""

	case C.VALUE_KIND_STRING:
		return C.GoStringN(v.string.data, C.int(v.string.len))

	case C.VALUE_KIND_BUFFER:
		b := C.GoBytes(
			unsafe.Pointer(v.string.data), C.int(v.string.len),
		)
		return string(b)
	}
}

func (v CValue) GetBytes() []byte {
	switch v.kind {
	default:
		return nil

	case C.VALUE_KIND_STRING:
		s := C.GoStringN(v.string.data, C.int(v.string.len))
		return []byte(s)

	case C.VALUE_KIND_BUFFER:
		return C.GoBytes(
			unsafe.Pointer(v.string.data), C.int(v.string.len),
		)
	}
}

func (v CValue) ArrayLen() int {
	if v.kind != C.VALUE_KIND_ARRAY {
		return 0
	}

	return int(v.array.len)
}

func (v CValue) ArrayForEach(fn func(i int, v *CValue)) {
	if v.kind != C.VALUE_KIND_ARRAY {
		return
	}

	n := int(v.array.len)
	for i := 0; i < n; i++ {
		vi := C.bridge_get_array_item(
			v.array.items,
			C.size_t(i),
		)

		fn(i, (*CValue)(vi))
	}
}

func (v CValue) GetPromise() *CPromise {
	return (*CPromise)(v.promise)
}

func (v CValue) PromiseResolve(val *CValue) {
	v.GetPromise().Resolve(val)
}

func (v CValue) PromiseReject(err *CValue) {
	v.GetPromise().Reject(err)
}

func (v *CValue) Free() error {
	kind := C.free_value((*C.struct_value_t)(v))

	// verify the value was cleared
	if kind != C.VALUE_KIND_UNDEFINED {
		return UnknownValueKindError{
			Kind: int(v.kind),
		}
	}

	return nil
}

func (v *CValue) Pointer() uintptr {
	return uintptr(unsafe.Pointer(v))
}
