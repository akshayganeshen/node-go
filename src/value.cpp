#include <stdio.h>

#include "./value.h"

using std::string;

using v8::Array;
using v8::Context;
using v8::FunctionCallbackInfo;
using v8::Isolate;
using v8::Local;
using v8::Number;
using v8::String;
using v8::Value;

static struct value_t *AllocValue(
  Isolate *isolate,
  const Local<Value> &val
);

static struct value_t *AllocBoolean(
  const Local<Value> &val
) {
  const unsigned char bol = val->IsTrue() ? 1 : 0;
  return alloc_value_boolean(bol);
}

static struct value_t *AllocNumber(
  const Local<Value> &val
) {
  Local<Number> num = val.As<Number>();
  return alloc_value_number(num->Value());
}

static struct value_t *AllocString(
  Isolate *isolate,
  const Local<Value> &val
) {
  String::Utf8Value ustr(isolate, val);
  return alloc_value_string(*ustr, ustr.length());
}

static struct value_t *AllocArray(
  Isolate *isolate,
  const Local<Value> &val
) {
  Local<Context> context = isolate->GetCurrentContext();
  Local<Array> arr = val.As<Array>();

  size_t len = arr->Length();
  struct value_t **items = new struct value_t *[len];

  for (size_t i = 0; i < len; i++) {
    Local<Value> vi = arr->Get(context, i).ToLocalChecked();
    items[i] = AllocValue(isolate, vi);
  }

  struct value_t *res = alloc_value_array(items, len);

  delete []items;
  return res;
}

static struct value_t *AllocValue(
  Isolate *isolate,
  const Local<Value> &val
) {
  if (val->IsNull()) {
    return alloc_value_null();
  }

  if (val->IsUndefined()) {
    return alloc_value_undefined();
  }

  if (val->IsBoolean()) {
    return AllocBoolean(val);
  }

  if (val->IsNumber()) {
    return AllocNumber(val);
  }

  if (val->IsString()) {
    return AllocString(isolate, val);
  }

  if (val->IsArray()) {
    return AllocArray(isolate, val);
  }

  fprintf(stderr, "value: unknown v8 value type\n");
  return NULL;
}

size_t GetNumArgs(
  const FunctionCallbackInfo<Value> &cb
) {
  return cb.Length();
}

struct array_t *AllocArgs(
  const FunctionCallbackInfo<Value> &cb
) {
  const int n = cb.Length();
  if (n == 0) {
    return alloc_array(NULL, 0);
  }

  struct value_t **items = new struct value_t *[n];

  Isolate *isolate = cb.GetIsolate();
  for (int i = 0; i < n; i++) {
    Local<Value> v = cb[i];
    items[i] = AllocValue(isolate, v);
  }

  struct array_t *args = alloc_array(items, n);

  delete []items;
  return args;
}

void FreeArgs(struct array_t *args) {
  free_array(args);
}
