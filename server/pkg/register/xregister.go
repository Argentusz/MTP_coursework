package register

// Using pointers the way the following code does is absolutely ridiculous
// This file not only hacky but also unsafe and, probably, hardware-dependant
//
// The only reason for it to be that way is to fulfill idiotic and unnecessary ideas
// of its lunatic author who should have written everything in C but was willing to be "classy"

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"unsafe"
)

// XRegister is a physical register which is associated with three different "virtual" registers
// its "extended" (ext) vreg is 32 bit long and has operates with all the data a physical register has
// "high" and "low" vregs are 16 bit long and operate with upper or lower bits of register memory, without interfering each other
type XRegister struct {
	ext *types.Word32
	hig *types.Word16
	low *types.Word16
}

func InitXRegister() XRegister {
	var data types.Word32
	dataptr := unsafe.Pointer(&data)
	hptri := *(*uint64)(unsafe.Pointer(&dataptr)) + 2

	var reg XRegister
	reg.ext = (*types.Word32)(dataptr)
	reg.hig = (*types.Word16)(dataptr)
	reg.low = *(**types.Word16)(unsafe.Pointer(&hptri))

	return reg
}
