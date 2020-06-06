#ifndef __LIB_FUNCTION_H__
#define __LIB_FUNCTION_H__

#include "./promise.h"
#include "./value.h"

struct fn_state_t;
struct fn_context_t;

typedef void (*fn_return)(
  struct fn_context_t  *ctx,
  struct value_t       *v
);

typedef void (*fn_throw)(
  struct fn_context_t  *ctx,
  struct value_t       *v
);

typedef struct promise_t *(*fn_new_promise)(
  struct fn_context_t  *ctx
);

// NOTE: Promise has it's own delete function embedded inside.

struct fn_callbacks_t {
  fn_return               return_func;
  fn_throw                throw_func;
  fn_new_promise          new_promise_func;
};

struct fn_context_t {
  struct fn_state_t      *state;
  struct fn_callbacks_t   callbacks;

  struct array_t         *args;
};

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

struct fn_callbacks_t function_callbacks(
  struct fn_context_t  *ctx
);
size_t function_num_arguments(
  struct fn_context_t  *ctx
);
struct value_t *function_argument(
  struct fn_context_t  *ctx,
  size_t                i
);

void function_return(
  struct fn_context_t  *ctx,
  struct value_t       *val
);
void function_throw(
  struct fn_context_t  *ctx,
  struct value_t       *err
);

// WARN: Function context may be NULL after doing async work.
struct promise_t *function_new_promise(
  struct fn_context_t  *ctx
);

// NOTE: Function context may be NULL, but the callback will still run.
struct promise_t *function_callback_new_promise(
  struct fn_callbacks_t cbs,
  struct fn_context_t  *ctx
);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __LIB_FUNCTION_H__ */
