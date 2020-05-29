#include <node.h>
#include <stdio.h>
// include the header file generated from the Go build
#include "../lib/libgo.h"

void Initialize(Local<Object> exports) {
}

NODE_MODULE(NODE_GYP_MODULE_NAME, Initialize)
