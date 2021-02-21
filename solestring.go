package solestring

// Remove -DSOLESTRING_PACK to disable pointer tagging (slightly faster but less memory-efficient)

// #cgo CFLAGS: -std=c99 -g -O2 -Wall -Wpedantic -Wno-unused-variable -Itidwall_hashmap -DSOLESTRING_PACK
// #include <stdlib.h>
// #include "solestring.h"
import "C"
import (
	"encoding/binary"
	"sync"
	"unsafe"
)

func NewStore() *Store {
	return &Store{
		hmap: C.hmap_new(),
	}
}

type Store struct {
	hmap *C.struct_hashmap
	mu   sync.RWMutex
}

func (o *Store) LoadOrStore(s string) (actual String, loaded bool) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	o.mu.RLock()
	p := C.hmap_get(o.hmap, cs)
	o.mu.RUnlock()

	if p == nil {
		o.mu.Lock()
		ok := C.hmap_put(o.hmap, cs)
		o.mu.Unlock()
		if !ok {
			return 0, false
		}

		o.mu.RLock()
		p = C.hmap_get(o.hmap, cs)
		o.mu.RUnlock()

		return String(unsafe.Pointer(p)), false
	}
	return String(unsafe.Pointer(p)), true
}

func (o *Store) Close() {
	C.hmap_free(o.hmap)
}

type String uintptr

func (s String) Value() string {
	// Pointer
	if s&1 == 0 {
		cs := (*C.char)(unsafe.Pointer(s))
		return C.GoString(cs)
	}

	// Tagged pointer
	bytes := make([]byte, unsafe.Sizeof(s))
	binary.LittleEndian.PutUint64(bytes, uint64(s>>8))
	var eos int
	for eos = 0; eos < len(bytes)-1; eos++ {
		if bytes[eos] == '\x00' {
			break
		}
	}
	return string(bytes[:eos])
}
