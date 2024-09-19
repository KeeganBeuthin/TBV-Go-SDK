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
		p := unsafe.Pointer(ptr)
		for i := 0; ; i++ {
			if *(*byte)(p) == 0 {
				length = int32(i)
				break
			}
			p = unsafe.Add(p, 1)
		}
	}
	return unsafe.String(ptr, int(length))
}

func PtrToString(ptr *byte) string {
	if ptr == nil {
		return ""
	}
	var bytes []byte
	for *ptr != 0 {
		bytes = append(bytes, *ptr)
		ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + 1))
	}
	return string(bytes)
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

// ReadHtmlCode processes the given HTML code
func ReadHtmlCode(htmlCode string) (string, error) {
	if htmlCode == "" {
		return "", fmt.Errorf("error: empty HTML code provided")
	}
	// Here you can add any processing logic if needed
	return htmlCode, nil
}

func CreateErrorResult(message string) *byte {
	errorMessage := fmt.Sprintf("Error: %s", message)
	return StringToPtr(errorMessage)
}

func CreateSuccessResult(message string) *byte {
	return StringToPtr(message)
}

//go:wasmimport env malloc
func wasmMalloc(size uint32) int32

func malloc(size int32) unsafe.Pointer {
	if size <= 0 {
		return nil
	}
	ptr := wasmMalloc(uint32(size))
	return unsafe.Pointer(uintptr(ptr))
}

// Free is a placeholder for freeing memory (if needed)
func Free(ptr unsafe.Pointer) {
	// Implementation depends on whether WebAssembly environment provides a free function
	// For now, we'll leave it empty
}
