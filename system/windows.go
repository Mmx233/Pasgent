package system

import (
	"fmt"
	"github.com/Mmx233/Pasgent/tools"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const (
	WmCopyData      = 0x004a
	AgentCopyDataId = 0x804e50ba
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procOpenFileMapping = kernel32.NewProc("OpenFileMappingW")
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procCreateWindowEx   = user32.NewProc("CreateWindowExW")
	procDefWindowProc    = user32.NewProc("DefWindowProcW")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
	procGetMessage       = user32.NewProc("GetMessageW")
	procRegisterClassEx  = user32.NewProc("RegisterClassExW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
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

type COPYDATASTRUCT struct {
	DwData uintptr
	CbData uint32
	LpData uintptr
}

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     windows.Handle
	HIcon         windows.Handle
	HCursor       windows.Handle
	HbrBackground windows.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       windows.Handle
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func getModuleHandle() (windows.Handle, error) {
	var hInstance windows.Handle
	err := windows.GetModuleHandleEx(0, nil, &hInstance)
	if err != nil {
		return 0, fmt.Errorf("failed to get module handle: %v", err)
	}
	return hInstance, nil
}

type WndProc func(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr

func CreateHiddenMessageWindow(className, windowName string, wndProc WndProc) error {
	_className, err := syscall.UTF16PtrFromString(className)
	if err != nil {
		return fmt.Errorf("failed to parse class name: %v", err)
	}
	_windowName, err := syscall.UTF16PtrFromString(windowName)
	if err != nil {
		return fmt.Errorf("failed to parse window name: %v", err)
	}

	hInstance, err := getModuleHandle()
	if err != nil {
		return err
	}

	wndClass := WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     hInstance,
		LpszClassName: _className,
	}
	if _, _, err := procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wndClass))); err != nil && err.(syscall.Errno) != 0 {
		return fmt.Errorf("failed to register windows message class: %v", err)
	}

	hwnd, _, err := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(_className)),
		uintptr(unsafe.Pointer(_windowName)),
		0,
		0, 0, 0, 0,
		0,
		0,
		uintptr(wndClass.HInstance),
		0,
	)
	if hwnd == 0 {
		return fmt.Errorf("failed to create Windows message class: %v", err)
	}

	var msg MSG
	for {
		ret, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
	return nil
}

type AgentRequestHandler func(data []byte) error

func NewAgentWndProc(handler AgentRequestHandler) WndProc {
	return func(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
		switch msg {
		case WmCopyData:
			cds := (*COPYDATASTRUCT)(unsafe.Pointer(lParam))
			if cds.DwData == AgentCopyDataId {
				err := handler(tools.PtrToSlice(cds.LpData, int(cds.CbData)))
				if err != nil {
					return 0
				}
				return 1
			}
		}
		ret, _, _ := procDefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}
}
