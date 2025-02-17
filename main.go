package main

// #cgo LDFLAGS: -lX11
import "C"

import (
	"github.com/dasciam/autoclicker-mcpe-go/x"
	"time"
)

func main() {
	display := x.Open()

	for {
		info := display.QueryPointer()
		wnd, _ := display.InputFocus()
		title := wnd.Title(display)

		if title == "Minecraft" && info.Mask&(1<<8) > 0 {
			attr := wnd.Attributes(display)

			wWidth := attr.Width / 2
			if wWidth+5 > info.RootX && wWidth-5 < info.RootX {
				display.SendMouseEvent(wnd)
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}
