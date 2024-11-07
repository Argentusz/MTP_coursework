package main

import (
	"testing"
	"unsafe"
)

func TestRandomStuff(t *testing.T) {
	uval1 := 0x3fc00000
	uval2 := 0x404ccccd

	fval1 := *((*float32)(unsafe.Pointer(&uval1)))
	fval2 := *((*float32)(unsafe.Pointer(&uval2)))
	fres := fval1 + fval2

	ures := *((*uint32)(unsafe.Pointer(&fres)))

	t.Log(uval1, uval2, ures)
	t.Log(fval1, fval2, fres)
}
