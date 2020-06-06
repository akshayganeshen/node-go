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
