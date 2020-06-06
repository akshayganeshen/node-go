# Node Go

Go module to easily create Node.js Native Addons.

## Developing

### CGo

Add shared headers in `lib/include/`. These headers can declare functions and
`extern` values, but should not include any definitions.

Implement any declared functions in a separate Go/C file under `lib/cgo/`.
Definitions cannot appear in a Go file that uses the `//export` directive. See
the cgo guide for info about this restriction.

### Go

WIP

### C++

Add addon code in `src/`.

Add an entry to `binding.gyp` to include it in the build.

## Running

An example set of handlers are included in `lib/exports.go`. The handlers are
loaded from `src/addon.cpp` to expose them as a Node module.

At a high-level, the handlers look like so:

```go
package main

/*
#ifndef __LIB_EXPORTS_GO__
#define __LIB_EXPORTS_GO__

#include "./include/function.h"

#endif // __LIB_EXPORTS_GO__
*/
import "C"

import (
  "strings"

  "node-go/lib/function"
  "node-go/lib/value"
)

func exAsyncHandler(ctx function.Context) error {
  result, err := ctx.Callbacks().NewPromise()
  if err != nil {
    return err
  }

  go func() {
    result.Resolve(value.String("goroutine resolve"))
  }()

  return ctx.Callbacks().Return(result.(*value.GoPromise))
}

func exDescribeArg(arg value.Js) string {
  if arg.IsNull() {
    return "null"
  }
  if arg.IsUndefined() {
    return "undefined"
  }

  if arg.IsBoolean() {
    return "boolean"
  }

  if arg.IsNumber() {
    return "number"
  }

  if arg.IsString() {
    return "string"
  }

  if arg.IsArray() {
    vs := arg.GetArray()
    items := make([]string, len(vs))
    for i, v := range vs {
      items[i] = exDescribeArg(v)
    }

    body := strings.Join(items, ", ")
    return "[ " + body + " ]"
  }

  return "other"
}

func exDescribeHandler(ctx function.Context) error {
  descs := make([]string, ctx.NumArguments())
  for i := range descs {
    descs[i] = exDescribeArg(ctx.GetArgument(i))
  }

  result := value.String("got args: " + strings.Join(descs, ", "))
  return ctx.Callbacks().Return(result)
}

var (
  wrappedExAsyncHandler    = WrapHandler(exAsyncHandler)
  wrappedExDescribeHandler = WrapHandler(exDescribeHandler)
)

//export ExAsyncHandler
func ExAsyncHandler(req *C.struct_fn_context_t) {
  wrappedExAsyncHandler(req)
}

//export ExDescribeHandler
func ExDescribeHandler(req *C.struct_fn_context_t) {
  wrappedExDescribeHandler(req)
}
```

The addon binds those handlers to `async` and  `describe`:

```js
const m = require('./build/Release/nodego');
console.log(m);
// { async: [Function: async], describe: [Function: describe] }
m.describe('hello', true, 1);
// 'got args: string, boolean, number'
m.describe('hello', true, 1, [ 'world', false, null ]);
// 'got args: string, boolean, number, [ string, boolean, null ]'
m.describe('hello', [ undefined, 'world', null, [ undefined ] ]);
// 'got args: string, [ undefined, string, null, [ undefined ] ]'

m.async().then((val) => console.log(`'${val}'`));
// Promise { <pending> }
// 'goroutine resolve'
```
