package main

import (
	"github.com/dasciam/autoclicker-mcpe-go/interfaces"
	"github.com/dasciam/autoclicker-mcpe-go/platform"
	"runtime"
	"time"
)

func main() {
	display, err := interfaces.Open()
	if err != nil {
		panic(err)
	}

	for {
		pointer, err := display.Pointer()
		if err != nil {
			panic(err)
		}
		window, err := display.WindowFocus()
		if err != nil {
			panic(err)
		}
		title := window.Title()

		if title == "Minecraft" && pointer.LoadMask(platform.FlagLMB) {
			width, _ := window.Size()
			wWidth := width / 2
			if runtime.GOOS == "windows" || wWidth+5 > pointer.X && wWidth-5 < pointer.X {
				_ = window.EmulateMouseClick()
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}
