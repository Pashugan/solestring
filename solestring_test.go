package solestring

import (
	"testing"
)

func TestFetchSame(t *testing.T) {
	s1, s2 := "test", "test"
	if &s2 == &s1 {
		t.Error("strings must not point to the same memory")
	}
	o := NewStore()
	defer o.Close()
	p1, _ := o.LoadOrStore(s1)
	p2, _ := o.LoadOrStore(s2)
	if p1.Value() != s1 {
		t.Error("p1 value is wrong")
	}
	if p2.Value() != s2 {
		t.Error("p2 value is wrong")
	}
}

func TestFetchNotSame(t *testing.T) {
	s1, s2 := "test", "TEST123"
	if &s2 == &s1 {
		t.Error("strings must not point to the same memory")
	}
	o := NewStore()
	defer o.Close()
	p1, _ := o.LoadOrStore(s1)
	p2, _ := o.LoadOrStore(s2)
	if p2 == p1 {
		t.Error("results must not point to the same memory")
	}
	if p1.Value() != s1 {
		t.Error("p1 value is wrong")
	}
	if p2.Value() != s2 {
		t.Error("p2 value is wrong")
	}
}

func TestFound(t *testing.T) {
	charFrom, charTo := 'a', 'f'
	o := NewStore()
	defer o.Close()
	for i := charFrom; i <= charTo; i++ {
		for j := charFrom; j <= charTo; j++ {
			_, found := o.LoadOrStore(string(i) + string(j))
			if found {
				t.Error("string should not be found during initialisation")
			}
		}
	}
	for i := charFrom; i <= charTo; i++ {
		for j := charFrom; j <= charTo; j++ {
			_, found := o.LoadOrStore(string(i) + string(j))
			if !found {
				t.Error("string should be found after initialisation")
			}
		}
	}
	_, found := o.LoadOrStore(string(charTo) + string(charTo+1))
	if found {
		t.Error("this string should not be found after initialisation")
	}
}

func BenchmarkFetchSame(b *testing.B) {
	s := "test"
	o := NewStore()
	defer o.Close()
	o.LoadOrStore(s)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o.LoadOrStore("test")
	}
}
