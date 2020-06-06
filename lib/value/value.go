package value

// Kind represents a Go value kind (i.e. value type).
// Kind is not necessarily in-sync with the C value kinds.
type Kind int

const (
	VALUE_KIND_NULL Kind = iota
	VALUE_KIND_UNDEFINED
	VALUE_KIND_ERROR
	VALUE_KIND_BOOLEAN
	VALUE_KIND_NUMBER
	VALUE_KIND_STRING
	VALUE_KIND_BUFFER
	VALUE_KIND_ARRAY
	VALUE_KIND_OBJECT
	VALUE_KIND_PROMISE
)

// Js is an abstraction for a JavaScript value.
type Js interface {
	IsNull() bool
	IsUndefined() bool
	IsNil() bool
	IsError() bool
	IsBoolean() bool
	IsNumber() bool
	IsString() bool
	IsBuffer() bool
	IsStringOrBuffer() bool
	IsArray() bool
	IsObject() bool
	IsPromise() bool

	GetBoolean() bool
	GetNumber() float64
	GetString() string
	GetBytes() []byte
	GetArray() []Js
	GetObject() map[string]Js
	GetPromise() Promise
}

type Go interface {
	// Go.Alloc allocates a copy of the value data to pass to C.
	// Go.Alloc expects the caller to also call C.Free.
	Alloc() (C, error)

	// Go.Kind returns the Go value kind.
	// Go.Kind can be used to perform type-assertions on the value.
	Kind() Kind
}

type C interface {
	// C.Free deallocates value data allocated by Go.Alloc.
	// C.Free should not be called on memory allocated by C.
	Free() error

	// C.Pointer returns a raw pointer to the C value data.
	// C.Pointer can be cast with unsafe to a *C.struct_value_t.
	// C.Pointer can be cast with unsafe to a *C.struct_promise_t.
	Pointer() uintptr
}
