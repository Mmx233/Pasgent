package system

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procOpenFileMapping = kernel32.NewProc("OpenFileMappingW")
)

func OpenFileMapping(desiredAccess uint32, inheritHandle uintptr, name string) (windows.Handle, error) {
	utf16Name, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	handle, _, err := procOpenFileMapping.Call(
		uintptr(desiredAccess),
		inheritHandle,
		uintptr(unsafe.Pointer(utf16Name)),
	)

	if handle == 0 {
		return 0, fmt.Errorf("open file mapping failed: %v", err)
	}
	return windows.Handle(handle), nil
}
