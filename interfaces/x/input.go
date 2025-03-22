//go:build linux

package x

//#include <stdint.h>
//void XQueryPointer(void*, uint64_t, void*, void*, void*, void*, void*, void*, void*);
//void XSendEvent(void*, uint64_t, int, int64_t, void*);
//void XGetInputFocus(void*, void*, void*);
import "C"

import "unsafe"

type Pointer struct {
	Root  Window
	Child Window

	RootX, RootY int32
	X, Y         int32

	Mask uint32
}

type xButtonEvent struct {
	unionType int32 // From union _XEvent.

	eventType     int32          // of event
	serial        uint64         // # of last request processed by server
	sendEvent     int            // bool, true if this came from a SendEvent request
	display       unsafe.Pointer // our *Display, Display the event was read from
	window        Window         // "event" window it is reported relative to
	root          Window         // root window that the event occurred on
	subWindow     Window         // child window
	time          uint64         // milliseconds
	x, y          int32          // pointer x, y coordinates in event window
	xRoot, yRoot  int32          // coordinates relative to root
	state, button uint32         // key or button mask | detail
	sameScreen    int            // bool

	_ [96]byte // Padding from union _XEvent.
}

func (d *Display) QueryPointer() Pointer {
	var ptr Pointer

	C.XQueryPointer(
		unsafe.Pointer(d),
		(*Screen)(unsafe.Add(d.screens, 128*d.defaultScreen)).root.Uint64(), // d.screens[d.defaultScreen]->root
		unsafe.Pointer(&ptr.Root),
		unsafe.Pointer(&ptr.Child),
		unsafe.Pointer(&ptr.RootX),
		unsafe.Pointer(&ptr.RootY),
		unsafe.Pointer(&ptr.X),
		unsafe.Pointer(&ptr.Y),
		unsafe.Pointer(&ptr.Mask),
	)
	return ptr
}

func (d *Display) SendMouseEvent(window Window) {
	var evt xButtonEvent

	evt.unionType = 4  // ButtonPress
	evt.button = 1     // Button1
	evt.sameScreen = 1 // true

	C.XSendEvent(
		unsafe.Pointer(d),
		window.Uint64(),
		C.int(1),
		C.int64_t(1<<2),
		unsafe.Pointer(&evt),
	)

	evt.unionType = 5 // ButtonRelease

	C.XSendEvent(
		unsafe.Pointer(d),
		window.Uint64(),
		C.int(1),
		C.int64_t(1<<3),
		unsafe.Pointer(&evt),
	)
}

func (d *Display) InputFocus() (window Window, revertTo int) {
	C.XGetInputFocus(
		unsafe.Pointer(d),
		unsafe.Pointer(&window),
		unsafe.Pointer(&revertTo),
	)
	return
}
