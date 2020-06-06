#include <stdio.h>

#include "./callbacks.h"
#include "./promise.h"

using namespace v8;

static void ReturnFunc(
  struct fn_context_t *ctx,
  struct value_t *v
) {
  if (!ctx->state) {
    // already returned
    fprintf(stderr, "return: invalid function state\n");
    return;
  }

  struct fn_state_t                 *state    = ctx->state;
  Isolate                           *isolate  = state->isolate;
  const FunctionCallbackInfo<Value> *cb       = state->cb;

  HandleScope scope(isolate);
  ReturnValue<Value> rv = cb->GetReturnValue();

  switch (v->kind) {
    case VALUE_KIND_NULL:
      rv.SetNull();
      break;

    case VALUE_KIND_UNDEFINED:
      rv.SetUndefined();
      break;

    case VALUE_KIND_STRING:
      rv.Set(
        String::NewFromUtf8(
          state->isolate,
          v->string->data,
          NewStringType::kNormal,
          v->string->len
        ).ToLocalChecked()
      );
      break;

    case VALUE_KIND_PROMISE:
      rv.Set(
        GetPromise(v->promise)
      );
      break;

    default:
      fprintf(
        stderr, "return: unknown value kind '%d'\n", v->kind
      );
      break;
  }
}

static void ThrowFunc(
  struct fn_context_t *ctx,
  struct value_t *v
) {
  fprintf(stderr, "throw: unimplemented\n");
}

static struct promise_t *NewPromiseFunc(
  struct fn_context_t *ctx
) {
  if (!ctx->state) {
    // already returned
    fprintf(stderr, "promise: invalid function state\n");
    return NULL;
  }

  struct fn_state_t *state    = ctx->state;
  Isolate           *isolate  = state->isolate;

  HandleScope scope(isolate);
  Local<Context> context = isolate->GetCurrentContext();

  Local<Promise::Resolver> resolver =
    Promise::Resolver::New(context).ToLocalChecked();

  return AllocPromise(isolate, resolver);
}

struct fn_state_t *AllocState(
  const FunctionCallbackInfo<Value> &cb
) {
  struct fn_state_t *state = new struct fn_state_t;
  state->isolate = cb.GetIsolate();

  // state->rv = cb.GetReturnValue();
  state->cb = &cb;

  fprintf(stderr, "[debug] function: allocate state (%p)\n", state);
  return state;
}

void FreeState(struct fn_state_t *state) {
  if (state) {
    fprintf(stderr, "[debug] function: delete state (%p)\n", state);
  }
  delete state;
}

struct fn_callbacks_t GetCallbacks() {
  struct fn_callbacks_t callbacks = {
    .return_func      = &ReturnFunc,
    .throw_func       = &ThrowFunc,
    .new_promise_func = &NewPromiseFunc,
  };

  return callbacks;
}
