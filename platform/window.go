package platform

type Window interface {
	Title() string

	Size() (width, height int32)

	EmulateMouseClick() error
}
