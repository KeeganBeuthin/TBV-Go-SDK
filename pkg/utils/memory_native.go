// //go:build !js || !wasm
// // +build !js !wasm

// package utils

// import "unsafe"

// type nativeAllocator struct{}

// func (n nativeAllocator) Malloc(size int32) unsafe.Pointer {
// 	buf := make([]byte, size)
// 	return unsafe.Pointer(&buf[0])
// }

// func getMemoryAllocator() MemoryAllocator {
// 	return nativeAllocator{}
// }
