//go:build windows

package interfaces

import (
	"github.com/dasciam/autoclicker-mcpe-go/interfaces/win"
	"github.com/dasciam/autoclicker-mcpe-go/platform"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
	"sync/atomic"
	"time"
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
	if err := w.handle.EmulateMouseDown(); err != nil {
		return err
	}
	return w.handle.EmulateMouseUP()
}

type winDisplay struct {
	point atomic.Pointer[types.POINT]

	leftMouseButton atomic.Bool
	lmbGuard        atomic.Bool

	pointerHidden atomic.Bool

	close chan struct{}
}

func (*winDisplay) WindowFocus() (platform.Window, error) {
	return winWindow{
		handle: win.ForegroundWindow(),
	}, nil
}

func (d *winDisplay) Pointer() (platform.Pointer, error) {
	var x, y int32

	if point := d.point.Load(); point != nil {
		x, y = point.X, point.Y
	}

	var flags uint32

	// We only set flags if the
	// mouse control is grabbed
	// by the game (hidden).
	if d.pointerHidden.Load() {
		if d.leftMouseButton.Load() {
			flags |= platform.FlagLMB
		}
	}

	return platform.Pointer{
		X:    x,
		Y:    y,
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

	if err := mouse.Install(nil, events); err != nil {
		panic(err)
		return
	}

	ticker := time.NewTicker(time.Millisecond)

	defer func() {
		_ = mouse.Uninstall()
		close(events)
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			cursor, err := win.CursorInfo()
			if err != nil {
				panic(err)
			}
			d.pointerHidden.Store(cursor.Flags == 0)
		case ev := <-events:
			switch ev.Message {
			case 0x0201: // WM_LBUTTONDOWN
				if !d.leftMouseButton.CompareAndSwap(false, true) {
					d.lmbGuard.Store(true)
				}
			case 0x0202: // WM_LBUTTONUP
				if d.lmbGuard.CompareAndSwap(true, false) {
					continue
				}
				d.leftMouseButton.Store(false)
			}
			d.point.Store(&ev.POINT)
		case <-d.close:
			return
		}
	}
}
