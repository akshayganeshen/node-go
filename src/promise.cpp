#include "./promise.h"

using namespace v8;

static void ResolvePromise(
  struct promise_t *p,
  struct value_t *v
) {
  if (!p->state || !p->state->pending) {
    // already fulfilled
    fprintf(stderr, "resolve: already fulfilled\n");
    return;
  }

  struct promise_state_t  *state  = p->state;
  struct async_state_t    *async  = state->async;

  async->fulfillment = {
    .result = kPromiseResolve,
    .value  = alloc_value_copy(v),
  };

  // unset async state to ensure it doesn't get de-allocated early
  state->pending  = false;
  state->async    = NULL;

  fprintf(stderr, "[debug] resolve: trigger async (%p)\n", async);
  uv_async_send(async->handle);
}

static void RejectPromise(
  struct promise_t *p,
  struct value_t *e
) {
  if (!p->state || !p->state->pending) {
    // already fulfilled
    fprintf(stderr, "reject: already fulfilled\n");
    return;
  }

  struct promise_state_t  *state  = p->state;
  struct async_state_t    *async  = state->async;

  async->fulfillment = {
    .result = kPromiseReject,
    .value  = alloc_value_copy(e),
  };

  // unset async state to ensure it doesn't get de-allocated early
  state->pending  = false;
  state->async    = NULL;

  fprintf(stderr, "[debug] resolve: trigger async (%p)\n", async);
  uv_async_send(async->handle);
}

static struct resolver_state_t BuildResolver(
  Isolate *isolate,
  Local<Promise::Resolver> resolver
) {
  return {
    .isolate  = isolate,
    .handle   = Persistent<Promise::Resolver>(isolate, resolver),
  };
}

static void FulfillPromiseHandle(uv_async_t           *handle);
static void FreeAsync           (struct async_state_t *async);

static struct async_state_t *AllocAsync(
  Isolate *isolate,
  Local<Promise::Resolver> resolver
) {
  struct async_state_t *async = new struct async_state_t;
  async->resolver     = BuildResolver(isolate, resolver);
  async->fulfillment  = { .result = kPromisePending };
  async->handle       = new uv_async_t;

  int c = uv_async_init(
    uv_default_loop(),
    async->handle,
    &FulfillPromiseHandle
  );

  if (c) {
    fprintf(stderr, "async: failed to initialize async handle: %d\n", c);
    FreeAsync(async);
    return NULL;
  }

  // finalize async state to ensure it can be accessed later
  async->handle->data = async;  // back-reference to self
  fprintf(stderr, "[debug] async: allocate async (%p)\n", async);
  return async;
}

static void FreeAsync(
  struct async_state_t *async
) {
  if (async) {
    fprintf(stderr, "[debug] async: delete async (%p)\n", async);
    async->resolver.handle.Reset();

    delete async->handle;
    async->handle = NULL;
  }

  delete async;
}

static void FreeAsyncHandle(uv_handle_t *vh) {
  uv_async_t *handle = reinterpret_cast<uv_async_t *>(vh);
  struct async_state_t *async =
    static_cast<struct async_state_t *>(handle->data);
  FreeAsync(async);
}

static void FulfillPromiseHandle(
  uv_async_t *handle
) {
  struct async_state_t *async =
    static_cast<struct async_state_t *>(handle->data);
  fprintf(stderr, "[debug] async: fulfill async (%p)\n", async);

  Isolate *isolate = async->resolver.isolate;
  HandleScope scope(isolate);

  if (!isolate->InContext()) {
    fprintf(stderr, "promise: no running context\n");
    return;
  }

  // XXX: Handle more than just string values
  struct value_t *v = async->fulfillment.value;

  if (v->kind != VALUE_KIND_STRING) {
    fprintf(stderr, "promise: unsupported value kind: '%d'\n", v->kind);
    return;
  }

  fprintf(
    stderr,
    "[debug] promise: fulfill with string value length %zu\n",
    v->string->len
  );

  Local<Context>            context   = isolate->GetCurrentContext();
  Local<Promise::Resolver>  resolver  = async->resolver.handle.Get(isolate);
  Local<Value>              value     = String::NewFromUtf8(
    isolate,
    v->string->data,
    NewStringType::kNormal,
    v->string->len
  ).ToLocalChecked();

  Maybe<bool> fulfilled = Nothing<bool>();
  switch (async->fulfillment.result) {
    default:
      fprintf(
        stderr,
        "promise: invalid promise fulfillment state '%d'\n",
        async->fulfillment.result
      );
      return;

    case kPromisePending:
      fprintf(stderr, "promise: no fulfillment result set\n");
      return;

    case kPromiseResolve:
      fulfilled = resolver->Resolve(context, value);
      break;

    case kPromiseReject:
      fulfilled = resolver->Reject(context, value);
      break;
  }

  if (fulfilled.IsNothing() || fulfilled.FromJust() == false) {
    fprintf(stderr, "promise: failed to fulfill promise\n");
    return;
  }

  free_value(v);  // free copied value
  uv_close(reinterpret_cast<uv_handle_t *>(handle), &FreeAsyncHandle);
}

struct promise_t *AllocPromise(
  Isolate *isolate,
  Local<Promise::Resolver> resolver
) {
  struct promise_state_t *state = new struct promise_state_t;
  state->pending  = true;
  state->resolver = BuildResolver(isolate, resolver);
  state->async    = AllocAsync(isolate, resolver);

  struct promise_t *r = new struct promise_t;
  r->state    =   state;
  r->resolve  =  &ResolvePromise;
  r->reject   =  &RejectPromise;
  r->free     =  &FreePromise;

  fprintf(stderr, "[debug] promise: allocate promise (%p)\n", r);
  return r;
}

Local<Promise> GetPromise(struct promise_t *promise) {
  // XXX: Promise can be NULL if it is fulfilled before it is returned
  struct promise_state_t *state = promise->state;

  Isolate *isolate = promise->state->resolver.isolate;
  EscapableHandleScope scope(isolate);

  Local<Promise::Resolver> resolver = state->resolver.handle.Get(isolate);
  fprintf(stderr, "[debug] promise: promise handle (%p)\n", promise);
  return scope.Escape(resolver->GetPromise());
}

void FreePromise(struct promise_t *promise) {
  if (promise) {
    fprintf(stderr, "[debug] promise: delete promise (%p)\n", promise);
    struct promise_state_t *state = promise->state;
    promise->state = NULL;

    if (state) {
      state->resolver.handle.Reset();

      // delete allocated async handle, if still pending
      FreeAsync(state->async);
      state->async = NULL;
    }

    delete state;
  }

  delete promise;
}
