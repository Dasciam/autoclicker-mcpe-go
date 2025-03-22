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

type Pointer struct {
	X, Y int32
	Mask uint32
}

func (ptr Pointer) LoadMask(mask uint32) bool {
	return ptr.Mask&mask != 0
}
