#include "xpc_wrapper.h"
#include <dispatch/dispatch.h>
#include <xpc/connection.h>
#include <Block.h>
#include <stdio.h>

//
// types and errors are implemented as macros
// create some real objects to make them accessible to Go
//
xpc_type_t TYPE_ERROR = XPC_TYPE_ERROR;

xpc_type_t TYPE_ARRAY = XPC_TYPE_ARRAY;
xpc_type_t TYPE_DATA = XPC_TYPE_DATA;
xpc_type_t TYPE_DICT = XPC_TYPE_DICTIONARY;
xpc_type_t TYPE_INT64 = XPC_TYPE_INT64;
xpc_type_t TYPE_STRING = XPC_TYPE_STRING;
xpc_type_t TYPE_UUID = XPC_TYPE_UUID;

xpc_object_t ERROR_CONNECTION_INVALID = (xpc_object_t) XPC_ERROR_CONNECTION_INVALID;
xpc_object_t ERROR_CONNECTION_INTERRUPTED = (xpc_object_t) XPC_ERROR_CONNECTION_INTERRUPTED;
xpc_object_t ERROR_CONNECTION_TERMINATED = (xpc_object_t) XPC_ERROR_TERMINATION_IMMINENT;

//
// this is the Go event handler
//
extern void HandleXPCEvent(xpc_object_t, void *);

extern void DictSet(void *, const char *, xpc_object_t);
extern void ArraySet(void *, int, xpc_object_t);

//
// connect to Apple Blued service
//
xpc_connection_t XpcConnectBlued(void *event_handler) {
    char *service = "com.apple.blued";
    dispatch_queue_t queue = dispatch_queue_create(service, 0);
    xpc_connection_t conn = xpc_connection_create_mach_service(service, queue, XPC_CONNECTION_MACH_SERVICE_PRIVILEGED);

    // this is needed to make sure that the following block
    // has a reference to a valid memory location
    void *eh = *((void **)event_handler);

    xpc_connection_set_event_handler(conn,
        ^(xpc_object_t event) {
            HandleXPCEvent(event, (void *)&eh);
        }
    );

    xpc_connection_resume(conn);
    return conn;
}

void XpcSendMessage(xpc_connection_t conn, xpc_object_t message, bool release) {
    xpc_connection_send_message(conn,  message);
    xpc_connection_send_barrier(conn, ^{
        // Block is invoked on connection's target queue
        // when 'message' has been sent.
        puts("message delivered");
    });
    if (release) {
        xpc_release(message);
    }
}

void XpcArrayApply(void *v, xpc_object_t arr) {
  xpc_array_apply(arr, ^bool(size_t index, xpc_object_t value) {
    ArraySet(v, index, value);
    return true;
  });
}

void XpcDictApply(void *v, xpc_object_t dict) {
  xpc_dictionary_apply(dict, ^bool(const char *key, xpc_object_t value) {
    DictSet(v, key, value);
    return true;
  });
}

void XpcUUIDGetBytes(void *v, xpc_object_t uuid) {
   const uint8_t *src = xpc_uuid_get_bytes(uuid);
   uint8_t *dest = (uint8_t *)v;

   for (int i=0; i < sizeof(uuid_t); i++) {
     dest[i] = src[i];
   }
}
