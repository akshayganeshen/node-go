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
