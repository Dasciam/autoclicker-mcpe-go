//go:build linux

package x

//#include <stdint.h>
//#include <string.h>
//int XFetchName(void*, uint64_t, char**);
//void XGetWindowAttributes(void*, uint64_t, void*);
//void XFree(void*);
import "C"

import "unsafe"

type WindowAttributes struct {
	X, Y, Width, Height int32

	_ [120]byte
}

type Window uint64

func (w Window) Title(display *Display) string {
	var cName *C.char

	status := C.XFetchName(
		unsafe.Pointer(display),
		w.Uint64(),
		&cName,
	)
	if status == -1 {
		return ""
	}
	defer C.XFree(unsafe.Pointer(cName))

	return C.GoString(cName)
}

func (w Window) Attributes(display *Display) (attr WindowAttributes) {
	C.XGetWindowAttributes(
		unsafe.Pointer(display),
		w.Uint64(),
		unsafe.Pointer(&attr),
	)
	return
}

func (w Window) Uint64() C.uint64_t {
	return C.uint64_t(w)
}
