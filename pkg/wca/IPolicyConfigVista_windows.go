//go:build windows
// +build windows

package wca

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

func pcvSetDefaultEndpoint(pcv *IPolicyConfigVista, deviceID string, eRole ERole) (err error) {
	var ptr *uint16
	if ptr, err = syscall.UTF16PtrFromString(deviceID); err != nil {
		return
	}
	hr, _, _ := syscall.SyscallN(
		pcv.VTable().SetDefaultEndpoint,
		uintptr(unsafe.Pointer(pcv)),
		uintptr(unsafe.Pointer(ptr)),
		uintptr(uint32(eRole)))
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
