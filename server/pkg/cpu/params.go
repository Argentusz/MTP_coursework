package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) get(bytes byte) types.Word32 {
	read := *cpu.RRAM.SYS.TMP & ((1 << bytes) - 1)
	*cpu.RRAM.SYS.TMP >>= bytes
	return read
}

func (cpu *CPU) getOperator() types.Word32 {
	return cpu.get(types.OperatorSize)
}

func (cpu *CPU) getReg() types.Word32 {
	return cpu.get(types.RegTypeSize)
}

func (cpu *CPU) getFlag() types.Word32 {
	return cpu.get(types.FlagTypeSize)
}

func (cpu *CPU) getInt() types.Word32 {
	return cpu.get(types.IntTypeSize)
}

func (cpu *CPU) getAddr() types.Word32 {
	return cpu.get(types.AddressTypeSize)
}

func (cpu *CPU) getSrc() types.Word32 {
	return cpu.get(types.ValueSourceTypeSize)
}

func (cpu *CPU) getDest() types.Word32 {
	return cpu.get(types.ValueDestinationTypeSize)
}
