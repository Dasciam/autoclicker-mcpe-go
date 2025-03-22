//go:build windows

package interfaces

//#include <stdio.h>
//#include <stdlib.h>
//extern int CallNextHookEx(void*, int, void*, void*);
//int __process_mouse_hook(int code, void* w, void* l) {
//  exit(0);
//	return CallNextHookEx(0, code, w, l);
//}
import "C"

import (
	"github.com/dasciam/autoclicker-mcpe-go/interfaces/win"
	"github.com/dasciam/autoclicker-mcpe-go/platform"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
	"log"
	"sync/atomic"
)

type winWindow struct {
	handle win.WindowHandle
}

func (w winWindow) Title() string {
	return w.handle.Title()
}

func (w winWindow) Size() (width, height int32) {
	return w.handle.Size()
}

func (w winWindow) EmulateMouseClick() error {
	w.handle.EmulateMouseClick()
	return nil
}

type winDisplay struct {
	lastMouseEvent atomic.Pointer[types.MouseEvent]

	rmbPressed atomic.Bool
	guard      atomic.Bool

	close chan struct{}
}

func (*winDisplay) WindowFocus() (platform.Window, error) {
	return winWindow{
		handle: win.ForegroundWindow(),
	}, nil
}

func (d *winDisplay) Pointer() (platform.Pointer, error) {
	var point types.POINT

	if ev := d.lastMouseEvent.Load(); ev != nil {
		point = ev.POINT
	}

	var flags uint32

	if d.rmbPressed.Load() {
		flags |= platform.FlagLMB
	}

	return platform.Pointer{
		X:    point.X,
		Y:    point.Y,
		Mask: flags,
	}, nil
}

func (d *winDisplay) Close() error {
	close(d.close)
	return nil
}

func Open() (platform.Display, error) {
	display := &winDisplay{
		close: make(chan struct{}),
	}
	go display.startTicking()
	return display, nil
}

func (d *winDisplay) startTicking() {
	events := make(chan types.MouseEvent)

	if err := mouse.Install(func(c chan<- types.MouseEvent) types.HOOKPROC {
		return func(_ int32, wParam, _ uintptr) uintptr {
			switch wParam {
			case 0x0201: // WM_LBUTTONDOWN
				if !d.rmbPressed.CompareAndSwap(false, true) {
					d.guard.Store(true)
				}
				log.Println("Button is pressed")
			case 0x0202: // WM_LBUTTONUP
				if d.guard.CompareAndSwap(true, false) {
					return 0
				}
				d.rmbPressed.Store(false)
				log.Println("Button up")
			}
			return 0
		}
	}, events); err != nil {
		panic(err)
		return
	}

	defer func() {
		_ = mouse.Uninstall()
		close(events)
	}()

	for {
		select {
		case ev := <-events:
			d.lastMouseEvent.Store(&ev)
		case <-d.close:
			return
		}
	}
}
