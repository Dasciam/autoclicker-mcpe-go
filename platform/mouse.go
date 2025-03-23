package platform

const (
	_ = 1 << iota
	_
	_
	_
	_
	_
	_
	FlagLMB
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
