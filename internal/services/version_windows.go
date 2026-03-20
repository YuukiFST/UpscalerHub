package services

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modVersion                  = syscall.NewLazyDLL("version.dll")
	procGetFileVersionInfoSizeW = modVersion.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfoW     = modVersion.NewProc("GetFileVersionInfoW")
	procVerQueryValueW          = modVersion.NewProc("VerQueryValueW")
)

type vsFixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

// readPEVersion reads the embedded version resource from a Windows PE file (DLL/EXE).
func readPEVersion(path string) (string, error) {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", err
	}

	size, _, _ := procGetFileVersionInfoSizeW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0,
	)
	if size == 0 {
		return "", fmt.Errorf("GetFileVersionInfoSizeW returned 0")
	}

	buf := make([]byte, size)
	ret, _, _ := procGetFileVersionInfoW.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0,
		size,
		uintptr(unsafe.Pointer(&buf[0])),
	)
	if ret == 0 {
		return "", fmt.Errorf("GetFileVersionInfoW failed")
	}

	subBlock, _ := syscall.UTF16PtrFromString(`\`)
	var info *vsFixedFileInfo
	var infoLen uint32

	ret, _, _ = procVerQueryValueW.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(subBlock)),
		uintptr(unsafe.Pointer(&info)),
		uintptr(unsafe.Pointer(&infoLen)),
	)
	if ret == 0 || infoLen == 0 {
		return "", fmt.Errorf("VerQueryValueW failed")
	}

	major := info.FileVersionMS >> 16
	minor := info.FileVersionMS & 0xffff
	patch := info.FileVersionLS >> 16
	build := info.FileVersionLS & 0xffff

	return fmt.Sprintf("%d.%d.%d.%d", major, minor, patch, build), nil
}
