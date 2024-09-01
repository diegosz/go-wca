//go:build windows
// +build windows

package wca

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

func pvString(v int64) (s string) {
	if v == 0 {
		return
	}
	var us []uint16
	var i uint32
	var start = unsafe.Pointer(uintptr(v))
	for {
		u := *(*uint16)(unsafe.Pointer(uintptr(start) + 2*uintptr(i)))
		if u == 0 {
			break
		}
		us = append(us, u)
		i++
	}
	s = syscall.UTF16ToString(us)
	ole.CoTaskMemFree(uintptr(v))
	return
}

func stringToPropVariant(s string) (v PROPVARIANT, err error) {
	// TODO: Is this correct?
	// Maybe we should use ole.SysAllocString instead of UTF16PtrFromString
	v.VT = ole.VT_LPWSTR
	ptr, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		return v, err
	}
	v.Val = int64(uintptr(unsafe.Pointer(ptr)))
	return
}
