package register

type FlagsRegister byte

// TODO: Too lazy to template code for each flag rn, will be written on demand

func (f *FlagsRegister) FZ() bool {
	return (*f & 0b1) == 1
}

func (f *FlagsRegister) ZeroOn() {
	*f |= 0b1
}

func (f *FlagsRegister) ZeroOff() {
	*f &= 0b0
}

func (f *FlagsRegister) MovZero(value bool) {
	switch value {
	case true:
		f.ZeroOn()
	case false:
		f.ZeroOff()
	}
}

func (f *FlagsRegister) FC() FlagsRegister {
	return *f & 0b10
}

func (f *FlagsRegister) OverflowOn() {
	*f |= 0b1000
}
