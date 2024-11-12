package register

import "github.com/Argentusz/MTP_coursework/pkg/types"

type FlagsRegister byte

// FZ - Result is zero
// FC - Carry (unsigned operation overflow)
// FS - Results highest bit
// FO - Overflow (signed operation)
// FI - Interrupts enabled (switches only by user)
// FT - Step by step mode (switches only by supervisor)
// FU - Sudo mode (switches only by supervisor)
// TODO: Too lazy to template code for each flag rn, will be written on demand

func (f *FlagsRegister) F(ID types.Word32) bool {
	flagMask := FlagsRegister(1 << ID)
	return (*f & flagMask) == flagMask
}

func (f *FlagsRegister) FZ() bool {
	return (*f & 0b1) == 0b1
}

func (f *FlagsRegister) FZOn() {
	*f |= 0b1
}

func (f *FlagsRegister) FC() bool {
	return (*f & 0b10) == 0b10
}

func (f *FlagsRegister) FCOn() {
	*f |= 0b10
}

func (f *FlagsRegister) FS() bool {
	return (*f & 0b100) == 0b100
}

func (f *FlagsRegister) FSOn() {
	*f |= 0b100
}

func (f *FlagsRegister) FO() bool {
	return (*f & 0b1000) == 0b1000
}

func (f *FlagsRegister) FOOn() {
	*f |= 0b1000
}

func (f *FlagsRegister) FI() bool {
	return (*f & 0b10000) == 0b10000
}

func (f *FlagsRegister) FIOn() {
	*f |= 0b10000
}

func (f *FlagsRegister) FIOff() {
	*f &= 0b11101111
}

func (f *FlagsRegister) FT() bool {
	return (*f & 0b100000) == 0b100000
}

func (f *FlagsRegister) FTOn() {
	*f |= 0b100000
}

func (f *FlagsRegister) FTOff() {
	*f &= 0b11011111
}

func (f *FlagsRegister) FU() bool {
	return (*f & 0b1000000) == 0b1000000
}

func (f *FlagsRegister) FUOn() {
	*f |= 0b1000000
}

func (f *FlagsRegister) FUOff() {
	*f &= 0b10111111
}

func (f *FlagsRegister) Drop() {
	*f &= 0b11110000
}

func (f *FlagsRegister) OnUnsignedOperation(fz, fs, fc bool) {
	f.Drop()
	if fz {
		f.FZOn()
	}
	if fs {
		f.FCOn()
	}
	if fc {
		f.FOOn()
	}
}

func (f *FlagsRegister) OnSignedOperation(fz, fs, fo bool) {
	f.Drop()
	if fz {
		f.FZOn()
	}
	if fs {
		f.FSOn()
	}
	if fo {
		f.FOOn()
	}
}
