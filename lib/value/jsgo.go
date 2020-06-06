package value

// JsGo wraps a Go value to provide a Js interface.
type JsGo struct {
	Go
}

var _ Js = JsGo{}

func GoToJs(v Go) Js {
	return JsGo{
		Go: v,
	}
}

func (v JsGo) IsNull() bool {
	return v.Kind() == VALUE_KIND_NULL
}

func (v JsGo) IsUndefined() bool {
	return v.Kind() == VALUE_KIND_UNDEFINED
}

func (v JsGo) IsNil() bool {
	return v.IsNull() || v.IsUndefined()
}

func (v JsGo) IsError() bool {
	return v.Kind() == VALUE_KIND_ERROR
}

func (v JsGo) IsBoolean() bool {
	return v.Kind() == VALUE_KIND_BOOLEAN
}

func (v JsGo) IsNumber() bool {
	return v.Kind() == VALUE_KIND_NUMBER
}

func (v JsGo) IsString() bool {
	return v.Kind() == VALUE_KIND_STRING
}

func (v JsGo) IsBuffer() bool {
	return v.Kind() == VALUE_KIND_BUFFER
}

func (v JsGo) IsStringOrBuffer() bool {
	return v.IsString() || v.IsBuffer()
}

func (v JsGo) IsArray() bool {
	return v.Kind() == VALUE_KIND_ARRAY
}

func (v JsGo) IsObject() bool {
	return v.Kind() == VALUE_KIND_OBJECT
}

func (v JsGo) IsPromise() bool {
	return v.Kind() == VALUE_KIND_PROMISE
}

func (v JsGo) GetBoolean() bool {
	switch v.Kind() {
	default:
		return true

	case VALUE_KIND_NULL:
		return false
	case VALUE_KIND_UNDEFINED:
		return false

	case VALUE_KIND_BOOLEAN:
		return bool(v.Go.(Boolean))

	case VALUE_KIND_NUMBER:
		return float64(v.Go.(Number)) != 0

	case VALUE_KIND_STRING:
		return len(string(v.Go.(String))) != 0
	}
}

func (v JsGo) GetNumber() float64 {
	// NOTE: Converting string to number is not supported
	if v.Kind() != VALUE_KIND_NUMBER {
		return 0
	}

	return float64(v.Go.(Number))
}

func (v JsGo) GetString() string {
	switch v.Kind() {
	default:
		return ""

	case VALUE_KIND_STRING:
		return string(v.Go.(String))

	case VALUE_KIND_BUFFER:
		return string([]byte(v.Go.(Buffer)))
	}
}

func (v JsGo) GetBytes() []byte {
	switch v.Kind() {
	default:
		return nil

	case VALUE_KIND_STRING:
		return []byte(string(v.Go.(String)))

	case VALUE_KIND_BUFFER:
		return []byte(v.Go.(Buffer))
	}
}

func (v JsGo) GetArray() []Js {
	if v.Kind() != VALUE_KIND_ARRAY {
		return nil
	}

	gvs := v.Go.(Array)
	n := gvs.Len()
	if n == 0 {
		return nil
	}

	jvs := make([]Js, n)
	gvs.ForEach(func(i int, value Go) {
		jvs[i] = GoToJs(value)
	})

	return jvs
}

func (v JsGo) GetObject() map[string]Js {
	if v.Kind() != VALUE_KIND_OBJECT {
		return nil
	}

	gmp := v.Go.(Object)
	n := gmp.Len()
	if n == 0 {
		return nil
	}

	jvs := make(map[string]Js, n)
	gmp.ForEach(func(key string, value Go) {
		jvs[key] = GoToJs(value)
	})

	return jvs
}

func (v JsGo) GetPromise() Promise {
	if v.Kind() != VALUE_KIND_PROMISE {
		return nil
	}

	return v.Go.(*GoPromise)
}
