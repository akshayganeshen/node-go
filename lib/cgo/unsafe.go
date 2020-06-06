package cgo

import (
	"unsafe"
)

// Unsafe is a uintptr helper to convert C types easily.
// Unsafe can be cast to a CValue, CPromise, or CFunctionContext.
type Unsafe uintptr

func UnsafeFrom(p unsafe.Pointer) Unsafe {
	return Unsafe(p)
}

func (u Unsafe) ToCValue() *CValue {
	return (*CValue)(unsafe.Pointer(u))
}

func (u Unsafe) ToCPromise() *CPromise {
	return (*CPromise)(unsafe.Pointer(u))
}

func (u Unsafe) ToCFunctionContext() *CFunctionContext {
	return (*CFunctionContext)(unsafe.Pointer(u))
}
