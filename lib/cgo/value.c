#ifndef __LIB_VALUE_C__
#define __LIB_VALUE_C__

#include <stdlib.h>
#include <string.h>

#include "../include/value.h"

struct value_t *alloc_value_null() {
  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (v) {
    v->kind = VALUE_KIND_NULL;
  }
  return v;
}

struct value_t *alloc_value_undefined() {
  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (v) {
    v->kind = VALUE_KIND_UNDEFINED;
  }
  return v;
}

struct value_t *alloc_value_boolean(unsigned char bol) {
  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (v) {
    v->kind     = VALUE_KIND_BOOLEAN;
    v->boolean  = bol;
  }
  return v;
}

struct value_t *alloc_value_number(double num) {
  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (v) {
    v->kind     = VALUE_KIND_NUMBER;
    v->number   = num;
  }
  return v;
}

struct value_t *alloc_value_string(const char *str, size_t len) {
  struct string_t *string = alloc_string(str, len);
  if (!string) {
    return NULL;
  }

  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (!v) {
    free_string(string);
    return NULL;
  }

  v->kind   = VALUE_KIND_STRING;
  v->string = string;
  return v;
}

struct value_t *alloc_value_buffer(void *buf, size_t len) {
  struct string_t *string = alloc_buffer(buf, len);
  if (!string) {
    return NULL;
  }

  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (!v) {
    free_string(string);
    return NULL;
  }

  v->kind   = VALUE_KIND_BUFFER;
  v->string = string;
  return v;
}

struct value_t *alloc_value_array(
  struct value_t **items,
  size_t len
) {
  struct array_t *array = alloc_array(items, len);
  if (!array) {
    return NULL;
  }

  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (!v) {
    free_array(array);
    return NULL;
  }

  v->kind   = VALUE_KIND_ARRAY;
  v->array  = array;
  return v;
}

struct value_t *alloc_value_object(
  struct array_t *keys,
  struct array_t *values,
  size_t len
) {
  struct object_t *object = alloc_object(keys, values, len);
  if (!object) {
    return NULL;
  }

  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (!v) {
    free_object(object);
    return NULL;
  }

  v->kind   = VALUE_KIND_OBJECT;
  v->object = object;
  return v;
}

struct value_t *alloc_value_promise(
  struct promise_t *promise
) {
  struct value_t *v =
    (struct value_t *)calloc(1, sizeof(struct value_t));
  if (v) {
    // as the header says, memory not managed
    v->kind     = VALUE_KIND_PROMISE;
    v->promise  = promise;
  }
  return v;
}

struct value_t *alloc_value_copy(struct value_t *val) {
  if (val) {
    switch (val->kind) {
    case VALUE_KIND_NULL:
      return alloc_value_null();

    case VALUE_KIND_UNDEFINED:
      return alloc_value_undefined();

    case VALUE_KIND_BOOLEAN:
      return alloc_value_boolean(val->boolean);

    case VALUE_KIND_NUMBER:
      return alloc_value_number(val->number);

    case VALUE_KIND_STRING:
      return alloc_value_string(val->string->data, val->string->len);

    case VALUE_KIND_BUFFER:
      return alloc_value_buffer(val->string->data, val->string->len);

    case VALUE_KIND_ARRAY:
      return alloc_value_array(val->array->items, val->array->len);

    case VALUE_KIND_OBJECT:
      return alloc_value_object(
        val->object->keys, val->object->values, val->object->len
      );

    case VALUE_KIND_PROMISE:
      return alloc_value_promise(val->promise);

    default:
      fprintf(
        stderr, "value: copy: unknown value kind '%d'\n", val->kind
      );
      break;
    }
  }

  return NULL;
}

int free_value(struct value_t *val) {
  if (val) {
    switch (val->kind) {
    case VALUE_KIND_NULL:
    case VALUE_KIND_UNDEFINED:
    case VALUE_KIND_BOOLEAN:
    case VALUE_KIND_NUMBER:
      break;

    case VALUE_KIND_PROMISE:
      // as the header says, memory not managed
      val->promise = NULL;
      break;

    case VALUE_KIND_STRING:
    case VALUE_KIND_BUFFER:
      free_string(val->string);
      val->string = NULL;
      break;

    case VALUE_KIND_ARRAY:
      free_array(val->array);
      val->array = NULL;
      break;

    default:
      fprintf(
        stderr, "value: free: unknown value kind '%d'\n", val->kind
      );
      // don't fall through and free data
      return val->kind;
    }

    val->kind = VALUE_KIND_UNDEFINED;
  }

  free(val);
  return VALUE_KIND_UNDEFINED;
}

struct string_t *alloc_string(const char *str, size_t len) {
  char *data = NULL;
  if (len) {
    data = (char *)calloc(len, sizeof(char));
    memcpy(data, str, len * sizeof(char));
  }

  struct string_t *s =
    (struct string_t *)calloc(1, sizeof(struct string_t));
  if (!s) {
    free(data);
    return NULL;
  }

  s->data = data;
  s->len  = len;
  return s;
}

struct string_t *alloc_buffer(void *buf, size_t len) {
  char *data = NULL;
  if (len) {
    data = (char *)calloc(len, sizeof(char));
    memcpy(data, buf, len * sizeof(char));
  }

  struct string_t *s =
    (struct string_t *)calloc(1, sizeof(struct string_t));
  if (!s) {
    free(data);
    return NULL;
  }

  s->data = data;
  s->len  = len;
  return s;
}

void free_string(struct string_t *str) {
  if (str) {
    free(str->data);
    str->data = NULL;
    str->len  = 0;
  }

  free(str);
}

struct array_t *alloc_array(
  struct value_t **its,
  size_t len
) {
  struct value_t **items = NULL;
  if (len && its) {
    items =
      (struct value_t **)calloc(len, sizeof(struct value_t *));
    memcpy(items, its, len * sizeof(struct value_t *));
  }

  struct array_t *arr =
    (struct array_t *)calloc(1, sizeof(struct array_t));
  if (!arr) {
    free(items);
    return NULL;
  }

  arr->items  = items;
  arr->len    = len;
  return arr;
}

void free_array(struct array_t *arr) {
  if (arr) {
    if (arr->items) {
      for (size_t i = 0; i < arr->len; i++) {
        free_value(arr->items[i]);
        arr->items[i] = NULL;
      }
    }

    free(arr->items);
    arr->items  = NULL;
    arr->len    = 0;
  }

  free(arr);
}

struct object_t *alloc_object(
  struct array_t *keys,
  struct array_t *values,
  size_t len
) {
  fprintf(stderr, "value: object: alloc: unimplemented\n");
  // unimplemented
  return NULL;
}

void free_object(struct object_t *obj) {
  fprintf(stderr, "value: object: free: unimplemented\n");
  // unimplemented
  free(obj);
}

#endif /* __LIB_VALUE_C__ */
