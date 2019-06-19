package api

import (
	"testing"
)

type to struct {
	F1 string
	F2 string
	f3 string
	f4 string
}

func (t *to) SetF3(f3 string) {
	t.f3 = f3
}
func (t *to) SetF4(f4 string) {
	t.f4 = f4
}

type from struct {
	F1 string
	f2 string
	F3 string
	f4 string
}

func (f *from) F2() string {
	return f.f2
}
func (f *from) F4() string {
	return f.f4
}

func TestCopyField(test *testing.T) {
	f := &from{
		F1: "f1",
		f2: "f2",
		F3: "f3",
	}
	t := new(to)
	err := copyField(t, f, []string{})
	if err != nil {
		test.Error(err)
	}
	if t.F1 != f.F1 {
		test.Fatal("direct copy not match")
	}
	if t.F2 != f.f2 {
		test.Fatal("copy from func not match")
	}
	if t.f3 != f.F3 {
		test.Fatal("copy to func not match")
	}
	if t.f4 != f.f4 {
		test.Fatal("func to func not match")
	}
}
