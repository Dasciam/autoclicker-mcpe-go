package platform

type Display interface {
	WindowFocus() (Window, error)
	Pointer() (Pointer, error)
}
