//go:build windows

package win

import (
	"syscall"
)

//goland:noinspection GoSnakeCaseUsage
const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")

	procMouseEvent    = user32.NewProc("mouse_event")
	procGetCursorInfo = user32.NewProc("GetCursorInfo")
)
