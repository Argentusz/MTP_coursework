package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) castSrcToImm(src types.Word32) types.Word32 {
	mode := src & types.SourceModeMask
	switch mode {
	case types.SourceRegMode:
		val, _ := cpu.RRAM.GetValue(src & types.SourceInverseMask)
		return types.Word32(val)
	case types.SourceIntMode:
		return src & types.SourceInverseMask
	case types.SourceAddrMode:
		addr, _ := cpu.RRAM.GetValue(src & types.SourceInverseMask)
		cpu.fetchUsrData(types.Address(addr))
		return *cpu.RRAM.SYS.MBR
	default:
		// SIGILL
		return src
	}
}

func (cpu *CPU) castSrcSize(src types.Word32) byte {
	mode := src & types.SourceModeMask
	switch mode {
	case types.SourceRegMode:
		return cpu.RRAM.GetRegSize(src & types.SourceInverseMask)
	case types.SourceIntMode:
		return 16
	case types.SourceAddrMode:
		return 32
	default:
		// SIGILL
		return 0
	}
}

func (cpu *CPU) castDstToImm(dst types.Word32) types.Word32 {
	mode := dst & types.DestinationModeMask
	switch mode {
	case types.DestinationRegMode:
		val, _ := cpu.RRAM.GetValue(dst & types.DestinationInverseMask)
		return types.Word32(val)
	case types.DestinationAddrMode:
		addr, _ := cpu.RRAM.GetValue(dst & types.DestinationInverseMask)
		cpu.fetchUsrData(types.Address(addr))
		return *cpu.RRAM.SYS.MBR
	default:
		// SIGILL
		return dst
	}
}

func (cpu *CPU) castDstToModeAddr(dst types.Word32) (bool, types.Word32) {
	mode := dst & types.DestinationModeMask
	switch mode {
	case types.DestinationRegMode:
		return true, dst & types.DestinationInverseMask
	case types.DestinationAddrMode:
		return false, cpu.castDstToImm(dst & types.DestinationInverseMask)
	default:
		// SIGILL
		return false, 0
	}
}

func (cpu *CPU) castAddrToImm(src types.Word32) types.Word32 {
	addr, _ := cpu.RRAM.GetValue(src)
	cpu.fetchUsrData(types.Address(addr))
	return *cpu.RRAM.SYS.MBR
}

func castValueSign(val types.Value, size byte) types.SValue {
	var valMask = types.Value((1 << (size - 1)) - 1)
	var valPure = val & valMask

	if valPure != val {
		return -types.SValue(valPure)
	}
	return types.SValue(val)
}

func castValueUnsign(val types.SValue, size byte) types.Value {
	if val < 0 {
		return types.Value(-val) | (1 << (size - 1))
	}

	return types.Value(val)
}
