package gomsp

/*
#cgo linux CFLAGS: -I./include
#cgo linux LDFLAGS: -L./lib ./lib/libMsp.a ./lib/libMMsp.a
#include <stdlib.h>
#include "t_msp.h"
*/
import "C"
import "unsafe"

func MspAttach(id uint16) int {
	return int(C.tMsp_Attach((C.ushort)(id)))
}

func MspDetach() int {
	return int(C.tMsp_Detach())
}

func MspPut(msg []byte, dstid uint16) int {

	putMsg := C.CString(string(msg))
	defer C.free(unsafe.Pointer(putMsg))
	return int(C.tMsp_Put(putMsg, C.int(len(msg)), (C.ushort)(dstid)))
}

func MspGet(timeout int) (msg []byte, len int32, srcid uint16, ret int) {
	getmsg := (*C.char)(C.calloc(8192, 1))
	defer C.free(unsafe.Pointer(getmsg))
	ret = int(C.tMsp_Get(getmsg,
		(*C.int)(unsafe.Pointer(&len)),
		(*C.ushort)(unsafe.Pointer(&srcid)),
		C.int(timeout)))
	msg = []byte(C.GoString(getmsg))
	return
}
