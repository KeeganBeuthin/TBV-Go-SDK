//go:build js && wasm
// +build js,wasm

package utils

import (
	"fmt"
	"unsafe"
)

func GoString(ptr *byte, length int32) string {
	if ptr == nil {
		return ""
	}
	if length < 0 {
		// Find null terminator
		end := ptr
		for *end != 0 {
			end = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(end)) + 1))
		}
		length = int32(uintptr(unsafe.Pointer(end)) - uintptr(unsafe.Pointer(ptr)))
	}
	return string(unsafe.Slice(ptr, length))
}

func StringToPtr(s string) *byte {
	bytes := []byte(s)
	ptr := malloc(int32(len(bytes) + 1))
	if ptr == nil {
		fmt.Println("Error: Failed to allocate memory")
		return nil
	}
	copy(unsafe.Slice((*byte)(ptr), len(bytes)+1), append(bytes, 0))
	return (*byte)(ptr)
}

func CreateErrorResult(message string) *byte {
	errorMessage := fmt.Sprintf("Error: %s", message)
	return StringToPtr(errorMessage)
}

func CreateSuccessResult(message string) *byte {
	return StringToPtr(message)
}

//go:wasmimport env malloc
//export malloc
func malloc(size int32) unsafe.Pointer
