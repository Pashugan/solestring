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

func (o *Store) LoadOrStore(s string) (actual *String, loaded bool) {
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
			return nil, false
		}

		o.mu.RLock()
		p = C.hmap_get(o.hmap, cs)
		o.mu.RUnlock()

		return (*String)(p), false
	}
	return (*String)(p), true
}

func (o *Store) Close() {
	C.hmap_free(o.hmap)
}

type String C.char

func (s *String) Value() string {
	cs := (*C.char)(s)
	p := uintptr(unsafe.Pointer(cs))

	// Pointer
	if p&1 == 0 {
		return C.GoString(cs)
	}

	// Tagged pointer
	bytes := make([]byte, unsafe.Sizeof(p))
	binary.LittleEndian.PutUint64(bytes, uint64(p>>8))
	var eos int
	for eos = 0; eos < len(bytes)-1; eos++ {
		if bytes[eos] == '\x00' {
			break
		}
	}
	return string(bytes[:eos])
}
