#ifndef __LIB_FUNCTION_C__
#define __LIB_FUNCTION_C__

#include "../include/function.h"

struct fn_callbacks_t function_callbacks(
  struct fn_context_t *ctx
) {
  return ctx->callbacks;
}

size_t function_num_arguments(
  struct fn_context_t *ctx
) {
  if (ctx->args) {
    return ctx->args->len;
  }

  return 0;
}

struct value_t *function_argument(
  struct fn_context_t *ctx,
  size_t i
) {
  return ctx->args->items[i];
}

void function_return(
  struct fn_context_t *ctx,
  struct value_t *val
) {
  ctx->callbacks.return_func(ctx, val);
}

void function_throw(
  struct fn_context_t *ctx,
  struct value_t *err
) {
  ctx->callbacks.throw_func(ctx, err);
}

struct promise_t *function_new_promise(
  struct fn_context_t *ctx
) {
  return ctx->callbacks.new_promise_func(ctx);
}

struct promise_t *function_callback_new_promise(
  struct fn_callbacks_t cbs,
  struct fn_context_t *ctx
) {
  return cbs.new_promise_func(ctx);
}

#endif /* __LIB_FUNCTION_C__ */
