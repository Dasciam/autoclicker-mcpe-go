//go:build linux

package interfaces

import (
	"github.com/dasciam/autoclicker-mcpe-go/interfaces/x"
	"github.com/dasciam/autoclicker-mcpe-go/platform"
)

//#cgo LDFLAGS: -lX11
import "C"

type x11window struct {
	display *x.Display
	window  x.Window
}

func (x x11window) EmulateMouseClick() error {
	x.display.SendMouseEvent(x.window)
	return nil
}

func (x x11window) Title() string {
	return x.window.Title(x.display)
}

func (x x11window) Size() (width, height int32) {
	attrs := x.window.Attributes(x.display)
	return attrs.Width, attrs.Height
}

type x11display struct {
	display *x.Display
}

func (x x11display) WindowFocus() (platform.Window, error) {
	focus, _ := x.display.InputFocus()
	return x11window{x.display, focus}, nil
}

func (x x11display) Pointer() (platform.Pointer, error) {
	ptr := x.display.QueryPointer()

	var flags uint32

	if ptr.Mask&(1<<8) != 0 {
		flags |= platform.FlagLMB
	}

	return platform.Pointer{
		X:    ptr.X,
		Y:    ptr.Y,
		Mask: flags,
	}, nil
}

func (x x11display) Close() error {
	return x.display.Close()
}

func Open() (platform.Display, error) {
	return x11display{
		display: x.Open(),
	}, nil
}
