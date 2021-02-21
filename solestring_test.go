package solestring

import (
	"fmt"
	"testing"
)

func hexDump(s string) string {
	out := ""
	for i := 0; i < len(s); i++ {
		out += fmt.Sprintf("%02x ", s[i])
	}
	return out
}

func TestFetchSamePacked(t *testing.T) {
	s1, s2 := "test", "test"
	if &s2 == &s1 {
		t.Error("strings must not point to the same memory")
	}
	o := NewStore()
	defer o.Close()
	p1, _ := o.LoadOrStore(s1)
	p2, _ := o.LoadOrStore(s2)
	if p2 != p1 {
		t.Error("results must point to the same memory")
	}
	if p1.Value() != s1 {
		t.Errorf("p1 value is wrong: want '%v', got '%v' (hex: %s)", s1, p1.Value(), hexDump(p1.Value()))
	}
	if p2.Value() != s2 {
		t.Errorf("p2 value is wrong: want '%v', got '%v' (hex: %s)", s1, p2.Value(), hexDump(p2.Value()))
	}
}

func TestFetchNotSamePacked(t *testing.T) {
	s1, s2 := "test", "TEST"
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
		t.Errorf("p1 value is wrong: want '%v', got '%v' (hex: %s)", s1, p1.Value(), hexDump(p1.Value()))
	}
	if p2.Value() != s2 {
		t.Errorf("p2 value is wrong: want '%v', got '%v' (hex: %s)", s1, p2.Value(), hexDump(p2.Value()))
	}
}

func TestFetchSameNotPacked(t *testing.T) {
	s1, s2 := "test-big-enough-string", "test-big-enough-string"
	if &s2 == &s1 {
		t.Error("strings must not point to the same memory")
	}
	o := NewStore()
	defer o.Close()
	p1, _ := o.LoadOrStore(s1)
	p2, _ := o.LoadOrStore(s2)
	if p2 != p1 {
		t.Error("results must point to the same memory")
	}
	if p1.Value() != s1 {
		t.Errorf("p1 value is wrong: want '%v', got '%v' (hex: %s)", s1, p1.Value(), hexDump(p1.Value()))
	}
	if p2.Value() != s2 {
		t.Errorf("p2 value is wrong: want '%v', got '%v' (hex: %s)", s1, p2.Value(), hexDump(p2.Value()))
	}
}

func TestFetchNotSameNotPacked(t *testing.T) {
	s1, s2 := "test-big-enough-string", "TEST-BIG-ENOUGH-STRING"
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
		t.Errorf("p1 value is wrong: want '%v', got '%v' (hex: %s)", s1, p1.Value(), hexDump(p1.Value()))
	}
	if p2.Value() != s2 {
		t.Errorf("p2 value is wrong: want '%v', got '%v' (hex: %s)", s1, p2.Value(), hexDump(p2.Value()))
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
