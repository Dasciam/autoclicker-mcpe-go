package platform

// Display represents the display and is responsible for receiving
// information about which window is active and information about
// the mouse pointer.
type Display interface {
	// WindowFocus returns the window that is in focus.
	// If a system call error occurred, it will be returned.
	WindowFocus() (Window, error)
	// Pointer returns information about the mouse pointer.
	// Pointer function returns the mouse position and flags
	// reflecting the current mouse state.
	Pointer() (Pointer, error)
}
