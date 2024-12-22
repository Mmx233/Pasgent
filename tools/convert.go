package tools

import "unsafe"

func PtrToSlice(addr uintptr, size int) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(addr)), size)
}
