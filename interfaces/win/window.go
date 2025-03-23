//go:build windows

package win

import (
	"syscall"
	"unsafe"
)

type WindowHandle uintptr

func (w WindowHandle) Title() string {
	const nChars = 512
	buf := make([]uint16, nChars)
	ret, _, _ := procGetWindowTextW.Call(
		uintptr(w),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(nChars),
	)
	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func (w WindowHandle) Size() (int32, int32) {
	width, _, _ := procGetSystemMetrics.Call(uintptr(SM_CXSCREEN))
	height, _, _ := procGetSystemMetrics.Call(uintptr(SM_CYSCREEN))
	return int32(width), int32(height)
}

func (w WindowHandle) EmulateMouseDown() error {
	_, _, _ = procMouseEvent.Call(
		uintptr(0x0002),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
	)
	return nil
}

func (w WindowHandle) EmulateMouseUP() error {
	_, _, _ = procMouseEvent.Call(
		uintptr(0x0004),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(0),
	)
	return nil
}

func ForegroundWindow() WindowHandle {
	w, _, _ := procGetForegroundWindow.Call()
	return WindowHandle(w)
}
