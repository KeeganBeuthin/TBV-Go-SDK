package ffi

import "unsafe"

//go:wasm-module env
//export query_rdf_tbv_cli
func QueryRdfTbvCli(queryPtr *byte, queryLen int32) *byte

//go:wasm-module env
//export malloc
func Malloc(size int32) unsafe.Pointer
