package goid

import (
	"bytes"
	"reflect"
	"runtime"
	"strconv"
	"unsafe"
)

//go:nosplit
func getgptr() unsafe.Pointer

//go:nosplit
func getg() interface{}

var goidOffset uintptr = func() uintptr {
	g := getg()
	if f, ok := reflect.TypeOf(g).FieldByName("goid"); ok {
		return f.Offset
	}
	panic("get goid offset failed.")
}()

func Get() uint64 {
	g := getgptr()
	p := (*uint64)(unsafe.Pointer(uintptr((g)) + goidOffset))
	return *p
}

func goid() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
