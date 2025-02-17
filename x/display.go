package x

//#include <stdint.h>
//void XCloseDisplay(void*);
//void* XOpenDisplay(void*);
import "C"

import "unsafe"

type Screen struct {
	_    [16]byte
	root Window
}

type Display struct {
	_             [224]byte // offset to default_screen.
	defaultScreen int32
	_             [4]byte        // int nscreens;
	screens       unsafe.Pointer // Screen* screens;
}

func Open() *Display {
	return (*Display)(C.XOpenDisplay(unsafe.Pointer(nil)))
}

func (d *Display) Close() error {
	C.XCloseDisplay(unsafe.Pointer(d))
	return nil
}
