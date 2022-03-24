package main

import (
	"syscall"
	"unsafe"
)

func LPCWSTR(str string) uintptr {
	lpcwstr, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		panic(err)
	}
	return uintptr(unsafe.Pointer(lpcwstr))
}

const (
	HWND_BROADCAST = uintptr(0xFFFF) // https://github.com/lxn/win/blob/7a0e89e/user32.go#L289
)

const (
	// https://pub.dev/documentation/win32/latest/win32/WM_FONTCHANGE-constant.html
	WM_FONTCHANGE = 0x001D
)
