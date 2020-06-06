#ifndef __LIB_VALUE_H__
#define __LIB_VALUE_H__

#include <stdio.h>

#define VALUE_KIND_NULL       ( -1  )
#define VALUE_KIND_UNDEFINED  (  0  )
#define VALUE_KIND_ERROR      (  1  )
#define VALUE_KIND_BOOLEAN    (  2  )
#define VALUE_KIND_NUMBER     (  3  )
#define VALUE_KIND_STRING     (  4  )
#define VALUE_KIND_BUFFER     (  5  )
#define VALUE_KIND_ARRAY      ( 10  )
#define VALUE_KIND_OBJECT     ( 11  )
#define VALUE_KIND_PROMISE    ( 20  )

struct string_t;
struct error_t;
struct array_t;
struct object_t;
struct promise_t;

struct value_t {
  int kind;

  unsigned char     boolean;
  double            number;

  struct error_t   *error;
  struct string_t  *string;
  struct array_t   *array;
  struct object_t  *object;
  struct promise_t *promise;  // memory not managed
};

struct error_t {
  int               code;
  struct string_t  *reason;
};

struct string_t {
  char             *data;
  size_t            len;
};

struct array_t {
  struct value_t  **items;
  size_t            len;
};

struct object_t {
  struct array_t   *keys;
  struct array_t   *values;
  size_t            len;
};

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

struct value_t *alloc_value_null();
struct value_t *alloc_value_undefined();
struct value_t *alloc_value_boolean(unsigned char bol);
struct value_t *alloc_value_number(double num);
struct value_t *alloc_value_string(const char *str, size_t len);
struct value_t *alloc_value_buffer(void       *buf, size_t len);

struct value_t *alloc_value_array(
  struct value_t **items,
  size_t len
);
struct value_t *alloc_value_object(
  struct array_t *keys,
  struct array_t *values,
  size_t len
);
struct value_t *alloc_value_promise(
  struct promise_t *promise
);

struct value_t *alloc_value_copy(struct value_t *val);

int free_value(struct value_t *val);  // returns value kind on error

// Lower-level helpers

struct string_t  *alloc_string(const char *str, size_t len);
struct string_t  *alloc_buffer(void       *buf, size_t len);
void              free_string (struct string_t *str);

struct array_t *alloc_array(
  struct value_t **its,
  size_t len
);
void free_array(struct array_t *arr);

struct object_t *alloc_object(
  struct array_t *keys,
  struct array_t *values,
  size_t len
);
void free_object(struct object_t *obj);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __LIB_VALUE_H__ */
