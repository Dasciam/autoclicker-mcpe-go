package platform

const (
	FlagLMB = 1 << iota
)

// Pointer is the state of the mouse pointer.
type Pointer struct {
	X, Y int32
	Mask uint32
}

// LoadMask loads the mask.
func (ptr Pointer) LoadMask(mask uint32) bool {
	return ptr.Mask&mask != 0
}
