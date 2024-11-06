package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) skip() {
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) mov() {
	param1 := cpu.getDest()
	param2 := cpu.getSrc()

	srcImm := types.Value(cpu.castSrcToImm(param2))
	isRegister, dstAddr := cpu.castDstToModeAddr(param1)
	if isRegister {
		cpu.RRAM.PutValue(dstAddr, srcImm)
		return
	}

	cpu.postUsrData(types.Address(dstAddr), types.Word32(srcImm), cpu.castSrcSize(param2))
}

func (cpu *CPU) add() {
	cpu._umathRS(func(a, b types.Value) types.Value { return a + b })
}

func (cpu *CPU) adc() {
	cpu._umathRS(func(a, b types.Value) types.Value {
		if cpu.RRAM.SYS.FLG.FC() {
			return a + b + 1
		}
		return a + b
	})
}

func (cpu *CPU) sub() {
	cpu._umathRS(func(a, b types.Value) types.Value { return a - b })
}

func (cpu *CPU) sbb() {
	cpu._umathRS(func(a, b types.Value) types.Value {
		if cpu.RRAM.SYS.FLG.FC() {
			return a - b - 1
		}
		return a - b
	})
}

func (cpu *CPU) mul() {
	cpu._umathRS(func(a, b types.Value) types.Value { return a * b })
}

func (cpu *CPU) div() {
	cpu._umathRS(func(a, b types.Value) types.Value { return a / b })
}

func (cpu *CPU) iadd() {
	cpu._imathRS(func(a, b types.SValue) types.SValue { return a + b })
}

func (cpu *CPU) iadc() {
	cpu._imathRS(func(a, b types.SValue) types.SValue {
		if cpu.RRAM.SYS.FLG.FC() {
			return a + b + 1
		}
		return a + b
	})
}

func (cpu *CPU) isub() {
	cpu._imathRS(func(a, b types.SValue) types.SValue { return a - b })
}

func (cpu *CPU) isbb() {
	cpu._imathRS(func(a, b types.SValue) types.SValue {
		if cpu.RRAM.SYS.FLG.FC() {
			return a - b - 1
		}
		return a - b
	})
}

func (cpu *CPU) idiv() {
	cpu._imathRS(func(a, b types.SValue) types.SValue { return a / b })
}

func (cpu *CPU) imul() {
	cpu._imathRS(func(a, b types.SValue) types.SValue { return a * b })
}

func (cpu *CPU) shl() {
	cpu._bitRI(func(a, b types.Value) types.Value { return a << b })
}

func (cpu *CPU) shr() {
	cpu._bitRI(func(a, b types.Value) types.Value { return a >> b })
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
	cpu._bitRR(func(a, b types.Value) types.Value { return a & b })
}

func (cpu *CPU) or() {
	cpu._bitRR(func(a, b types.Value) types.Value { return a | b })
}

func (cpu *CPU) xor() {
	cpu._bitRR(func(a, b types.Value) types.Value { return a ^ b })
}

func (cpu *CPU) jmp() {
	addr := cpu.getAddr()
	*cpu.RRAM.SYS.IR = cpu.castAddrToImm(addr)
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) call() {
	addr := cpu.getAddr()
	*cpu.RRAM.SYS.IRB = *cpu.RRAM.SYS.IR
	*cpu.RRAM.SYS.IR = cpu.castAddrToImm(addr)
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) ret() {
	*cpu.RRAM.SYS.IR = *cpu.RRAM.SYS.IRB
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) ei() {
	cpu.RRAM.SYS.FLG.Drop()
	cpu.RRAM.SYS.FLG.FIOn()
}

func (cpu *CPU) di() {
	cpu.RRAM.SYS.FLG.Drop()
	cpu.RRAM.SYS.FLG.FIOff()
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

func (cpu *CPU) _bitRR(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	srcRegID := cpu.getReg()

	dstVal, err1 := cpu.RRAM.GetValue(dstRegID)
	srcVal, err2 := cpu.RRAM.GetValue(srcRegID)
	if err1 != nil || err2 != nil {
		// SIGILL
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := fn(dstVal, srcVal)
	overflow := cpu.RRAM.PutValue(dstRegID, res)
	cpu.RRAM.SYS.FLG.OnUnsignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)
}

func (cpu *CPU) _bitRI(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	immVal := cpu.getInt()
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := fn(dstVal, types.Value(immVal))
	overflow := cpu.RRAM.PutValue(dstRegID, res)
	cpu.RRAM.SYS.FLG.OnUnsignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)

}

func (cpu *CPU) _umathRS(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcVal := types.Value(cpu.castSrcToImm(source))
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := fn(dstVal, srcVal)
	overflow := cpu.RRAM.PutValue(dstRegID, res)
	cpu.RRAM.SYS.FLG.OnUnsignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)
}

func (cpu *CPU) _imathRS(fn func(a, b types.SValue) types.SValue) {
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcVal := types.Value(cpu.castSrcToImm(source))
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		// SIGILL
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	srcSize := cpu.castSrcSize(source)

	dstSVal := castValueSign(dstVal, dstSize)
	srcSVal := castValueSign(srcVal, srcSize)

	res := fn(dstSVal, srcSVal)
	overflow := cpu.RRAM.PutValue(dstRegID, castValueUnsign(res, dstSize))
	cpu.RRAM.SYS.FLG.OnSignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)
}
