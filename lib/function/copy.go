package function

import (
	"node-go/lib/cgo"
	"node-go/lib/value"
)

// CopyC creates a Go copy of a cgo.CFunctionContext.
// CopyC returns a Go pointer, so it is garbage collected automatically.
func CopyC(ctx *cgo.CFunctionContext) *GoContext {
	return &GoContext{
		ContextHandle: ctx,
		CallbacksCopy: ctx.GetCallbacks(),
		Arguments:     CopyArguments(ctx),
		Running:       true,
	}
}

func CopyArguments(ctx *cgo.CFunctionContext) []value.Js {
	n := ctx.NumArguments()
	if n == 0 {
		return nil
	}

	vs := make([]value.Js, n)
	for i := range vs {
		cv := ctx.GetArgument(i)
		gv := value.GoToJs(value.CopyC(cv))

		vs[i] = gv
	}

	return vs
}
