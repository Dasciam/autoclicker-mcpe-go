package win

import (
	"errors"
	"github.com/moutend/go-hook/pkg/types"
	"unsafe"
)

type Cursor struct {
	CbSize       uint32  // DWORD
	Flags        uint32  // DWORD
	CursorHandle uintptr // HCURSOR
	PtScreenPos  types.POINT
}

func CursorInfo() (info Cursor, err error) {
	info.CbSize = uint32(unsafe.Sizeof(info))

	v, _, _ := procGetCursorInfo.Call(
		uintptr(unsafe.Pointer(&info)),
	)
	if v == 0 {
		err = errors.New("GetCursorInfo: failed to get cursor info")
	}
	return
}
