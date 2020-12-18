package oncestring

// #include <stdlib.h>
// #include "solestring.h"
import "C"
import (
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

func (o *Store) LoadOrStore(s string) (*String, bool) {
	cs := C.CString(s)
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
		return (*String)(cs), false
	}
	C.free(unsafe.Pointer(cs))
	return (*String)(p), true
}

func (o *Store) Close() {
	C.hashmap_free(o.hmap)
}

type String C.char

func (s *String) Value() string {
	return C.GoString((*C.char)(s))
}
