#ifndef __LIB_PROMISE_H__
#define __LIB_PROMISE_H__

struct promise_t;
struct promise_state_t;
struct value_t;

typedef void (*promise_resolve_func)(
  struct promise_t *promise,
  struct value_t   *value
);

typedef void (*promise_reject_func)(
  struct promise_t *promise,
  struct value_t   *error
);

typedef void (*promise_free_func)(
  struct promise_t *promise
);

struct promise_t {
  struct promise_state_t *state;
  promise_resolve_func    resolve;
  promise_reject_func     reject;
  promise_free_func       free;
};

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

void promise_resolve(
  struct promise_t *promise,
  struct value_t   *value
);
void promise_reject(
  struct promise_t *promise,
  struct value_t   *error
);
void promise_free(
  struct promise_t *promise
);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __LIB_PROMISE_H__ */
