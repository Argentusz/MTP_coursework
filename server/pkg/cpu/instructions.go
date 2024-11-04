package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) skip() {

}

func (cpu *CPU) mov() {
	//fmt.Println("Executing MOV")
	param1 := cpu.getDest()
	param2 := cpu.getSrc()

	//fmt.Printf("Encrypted params: %07b %018b\n", param1, param2)
	srcImm := types.Value(cpu.castSrcToImm(param2))
	isRegister, dstAddr := cpu.castDstToModeAddr(param1)
	//fmt.Printf("Decypher params: source immediate value = %d, destination address = %d, destination is register %v\n", srcImm, dstAddr, isRegister)
	if isRegister {
		cpu.RRAM.PutValue(dstAddr, srcImm)
		return
	}

	cpu.postUsrData(types.Address(dstAddr), types.Word32(srcImm), cpu.castSrcSize(param2))
}

func (cpu *CPU) add() {
	cpu._binaryRSop(func(a, b types.Value) types.Value { return a + b })
}

func (cpu *CPU) adc() {
	cpu._binaryRSop(func(a, b types.Value) types.Value { return a + b + types.Value(cpu.RRAM.SYS.FLG.FC()) })
}

func (cpu *CPU) sub() {
	cpu._binaryRSop(func(a, b types.Value) types.Value { return a - b })
}

func (cpu *CPU) sbb() {
	cpu._binaryRSop(func(a, b types.Value) types.Value { return a - b - types.Value(cpu.RRAM.SYS.FLG.FC()) })
}

func (cpu *CPU) mul() {
	cpu._binaryRSop(func(a, b types.Value) types.Value { return a * b })
}

func (cpu *CPU) imul() {
	panic("TODO: imul is NOI.")
}

func (cpu *CPU) div() {
	cpu._binaryRSop(func(a, b types.Value) types.Value { return a / b })
}

func (cpu *CPU) shl() {
	cpu._binaryRIop(func(a, b types.Value) types.Value { return a << b })
}

func (cpu *CPU) shr() {
	cpu._binaryRIop(func(a, b types.Value) types.Value { return a >> b })
}

func (cpu *CPU) sar() {
	dstRegID := cpu.getReg()
	immVal := cpu.getInt()
	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	highBit := dstVal & (1 << (dstSize - 1))
	cpu.RRAM.PutValue(dstRegID, (dstVal>>immVal)|highBit)
}

func (cpu *CPU) and() {
	cpu._binaryRRop(func(a, b types.Value) types.Value { return a & b })
}

func (cpu *CPU) or() {
	cpu._binaryRRop(func(a, b types.Value) types.Value { return a | b })
}

func (cpu *CPU) xor() {
	cpu._binaryRRop(func(a, b types.Value) types.Value { return a ^ b })
}

func (cpu *CPU) not() {
	dstRegID := cpu.getReg()
	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	mask := types.Value((1 << dstSize) - 1)
	cpu.RRAM.PutValue(dstRegID, dstVal^mask)
}

func (cpu *CPU) _binaryRRop(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	srcRegID := cpu.getReg()

	dstVal, err1 := cpu.RRAM.GetValue(dstRegID)
	srcVal, err2 := cpu.RRAM.GetValue(srcRegID)
	if err1 != nil || err2 != nil {
		// SIGILL
		return
	}

	cpu.RRAM.PutValue(dstRegID, fn(dstVal, srcVal))
}

func (cpu *CPU) _binaryRIop(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	immVal := cpu.getInt()
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	cpu.RRAM.PutValue(dstRegID, fn(dstVal, types.Value(immVal)))
}

func (cpu *CPU) _binaryRSop(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcVal := types.Value(cpu.castSrcToImm(source))
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	cpu.RRAM.PutValue(dstRegID, fn(dstVal, srcVal))
}
