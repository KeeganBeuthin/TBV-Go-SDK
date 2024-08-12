//go:build js && wasm
// +build js,wasm

package utils

import "unsafe"

type wasmAllocator struct{}

//go:wasm-module env
//export malloc
func wasmMalloc(size int32) unsafe.Pointer

func (w wasmAllocator) Malloc(size int32) unsafe.Pointer {
	return wasmMalloc(size)
}

func getMemoryAllocator() MemoryAllocator {
	return wasmAllocator{}
}
