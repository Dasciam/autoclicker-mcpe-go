package platform

// Window is a system window, responsible for getting its
// name, size and click emulation.
type Window interface {
	// Title returns the window name.
	// If a system call error occurred, returns an empty string.
	Title() string
	// Size returns the window size,
	// If a system call error occurred, returns zeroes.
	Size() (width, height int32)
	// EmulateMouseClick emulates a click.
	// If a system call error occurred, it will be returned.
	EmulateMouseClick() error
}
