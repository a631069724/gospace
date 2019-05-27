package carddecrypt

/*
#cgo LDFLAGS: -L C:/mingw/x86_64-w64-mingw32/lib
#cgo LDFLAGS: ./libdecrypt.a
#include <stdlib.h>
#include <malloc.h>

int HuBeiDecrypt(int,char*[]);
*/
import "C"

import (
	"os"
	"unsafe"
)

func CardDecrypt() {
	arg := make([](*C.char), 0)
	var i int
	var v string
	for i, v = range os.Args {
		char := C.CString(v)
		defer C.free(unsafe.Pointer(char))
		strptr := (*C.char)(unsafe.Pointer(char))
		arg = append(arg, strptr)
	}
	C.HuBeiDecrypt(C.int(i), (**C.char)(unsafe.Pointer(&arg[0])))
}
