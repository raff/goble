package goble

/*
#include <xpc/xpc.h>
#include "xpc_wrapper.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	r "reflect"
	"unsafe"
)

//
// minimal XPC support required for BLE
//

// a dictionary of things
type dict map[string]interface{}

// an array of things
type array []interface{}

var (
	CONNECTION_INVALID     = errors.New("connection invalid")
	CONNECTION_INTERRUPTED = errors.New("connection interrupted")
	CONNECTION_TERMINATED  = errors.New("connection terminated")
)

type EventHandler func(event dict, err error)

// XpcConnect connects to the Blued service
// (with some trickery to pass around the event callback handler)
func XpcConnect(eh EventHandler) C.xpc_connection_t {
	return C.XpcConnectBlued((*(*unsafe.Pointer)(unsafe.Pointer(&eh))))
}

//export HandleXPCEvent
func HandleXPCEvent(event C.xpc_object_t, p unsafe.Pointer) {
	t := C.xpc_get_type(event)
	eh := (*(*EventHandler)(unsafe.Pointer(&p)))

	if t == C.TYPE_ERROR {
		if event == C.ERROR_CONNECTION_INVALID {
			// The client process on the other end of the connection has either
			// crashed or cancelled the connection. After receiving this error,
			// the connection is in an invalid state, and you do not need to
			// call xpc_connection_cancel(). Just tear down any associated state
			// here.
			log.Println("connection invalid")
			eh(nil, CONNECTION_INVALID)
		} else if event == C.ERROR_CONNECTION_INTERRUPTED {
			log.Println("connection interrupted")
			eh(nil, CONNECTION_INTERRUPTED)
		} else if event == C.ERROR_CONNECTION_TERMINATED {
			// Handle per-connection termination cleanup.
			log.Println("connection terminated")
			eh(nil, CONNECTION_TERMINATED)
		} else {
			log.Println("got some error", event)
			eh(nil, fmt.Errorf("%v", event))
		}
	} else {
		eh(xpcToGo(event).(dict), nil)
	}
}

// goToXpc converts a go object to an xpc object
func goToXpc(o interface{}) C.xpc_object_t {
	return valueToXpc(r.ValueOf(o))
}

// valueToXpc converts a go Value to an xpc object
//
// note that not all the types are supported, but only the subset required for Blued
func valueToXpc(val r.Value) C.xpc_object_t {
	if !val.IsValid() {
		return nil
	}

	var xv C.xpc_object_t

	switch val.Kind() {
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		xv = C.xpc_int64_create(C.int64_t(val.Int()))

	case r.Uint, r.Uint8, r.Uint16, r.Uint32:
		xv = C.xpc_int64_create(C.int64_t(val.Uint()))

	case r.String:
		xv = C.xpc_string_create(C.CString(val.String()))

	case r.Map:
		xv = C.xpc_dictionary_create(nil, nil, 0)
		for _, k := range val.MapKeys() {
			v := valueToXpc(val.MapIndex(k))
			C.xpc_dictionary_set_value(xv, C.CString(k.String()), v)
			if v != nil {
				C.xpc_release(v)
			}
		}

	case r.Array, r.Slice:
		xv = C.xpc_array_create(nil, 0)
		l := val.Len()

		for i := 0; i < l; i++ {
			v := valueToXpc(val.Index(i))
			C.xpc_array_append_value(xv, v)
			if v != nil {
				C.xpc_release(v)
			}
		}

	case r.Interface, r.Ptr:
		xv = valueToXpc(val.Elem())

	default:
		log.Fatalf("unsupported %#v", val.String())
	}

	return xv
}

//export ArraySet
func ArraySet(u unsafe.Pointer, i C.int, v C.xpc_object_t) {
	a := *(*array)(u)
	a[i] = xpcToGo(v)
}

//export DictSet
func DictSet(u unsafe.Pointer, k *C.char, v C.xpc_object_t) {
	d := *(*dict)(u)
	d[C.GoString(k)] = xpcToGo(v)
}

// xpcToGo converts an xpc object to a go object
//
// note that not all the types are supported, but only the subset required for Blued
func xpcToGo(v C.xpc_object_t) interface{} {
	t := C.xpc_get_type(v)

	switch t {
	case C.TYPE_ARRAY:
		a := make(array, C.int(C.xpc_array_get_count(v)))
		C.XpcArrayApply(unsafe.Pointer(&a), v)
		return a

	case C.TYPE_DATA:
		return C.GoBytes(C.xpc_data_get_bytes_ptr(v), C.int(C.xpc_data_get_length(v)))

	case C.TYPE_DICT:
		d := make(dict)
		C.XpcDictApply(unsafe.Pointer(&d), v)
		return d

	case C.TYPE_INT64:
		return int64(C.xpc_int64_get_value(v))

	case C.TYPE_STRING:
		return C.GoString(C.xpc_string_get_string_ptr(v))

	default:
		log.Fatalf("unexpected type %#v, value %#v", t, v)
	}

	return nil
}
