#ifndef __LIB_PROMISE_C__
#define __LIB_PROMISE_C__

#include "../include/promise.h"

void promise_resolve(
  struct promise_t *promise,
  struct value_t   *value
) {
  promise->resolve(promise, value);
}

void promise_reject(
  struct promise_t *promise,
  struct value_t   *error
) {
  promise->reject(promise, error);
}

void promise_free(
  struct promise_t *promise
) {
  promise->free(promise);
}

#endif /* __LIB_PROMISE_C__ */
