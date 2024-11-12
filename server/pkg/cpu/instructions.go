package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"unsafe"
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
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcVal := types.Value(cpu.castSrcToImm(source))
	if srcVal == 0 {
		cpu.SIGFPE()
		return
	}

	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		cpu.SIGILL()
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := dstVal / srcVal
	overflow := cpu.RRAM.PutValue(dstRegID, res)
	cpu.RRAM.SYS.FLG.OnUnsignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)
}

func (cpu *CPU) rmd() {
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcVal := types.Value(cpu.castSrcToImm(source))
	if srcVal == 0 {
		cpu.SIGFPE()
		return
	}

	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		cpu.SIGILL()
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := dstVal % srcVal
	overflow := cpu.RRAM.PutValue(dstRegID, res)
	cpu.RRAM.SYS.FLG.OnUnsignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)
}

func (cpu *CPU) imov() {
	dst := cpu.getDest()
	src := cpu.getSrc()

	srcImm := types.Value(cpu.castSrcToImm(src))
	isRegister, dstAddr := cpu.castDstToModeAddr(dst)

	if isRegister {
		srcSize := cpu.castSrcSize(src)
		srcHighbit := srcImm >> (srcSize - 1)
		dstSize := cpu.RRAM.GetRegSize(dstAddr)
		srcImm &= (1 << (srcSize - 1)) - 1
		srcImm |= srcHighbit << (dstSize - 1)
		cpu.RRAM.PutValue(dstAddr, srcImm)
		return
	}

	cpu.postUsrData(types.Address(dstAddr), types.Word32(srcImm), cpu.castSrcSize(src))

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
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcVal := types.Value(cpu.castSrcToImm(source))
	dstVal, err := cpu.RRAM.GetValue(dstRegID)
	if err != nil {
		cpu.SIGILL()
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	srcSize := cpu.castSrcSize(source)

	dstSVal := castValueSign(dstVal, dstSize)
	srcSVal := castValueSign(srcVal, srcSize)
	if srcSVal == 0 {
		cpu.SIGFPE()
		return
	}

	res := dstSVal / srcSVal
	overflow := cpu.RRAM.PutValue(dstRegID, castValueUnsign(res, dstSize))
	cpu.RRAM.SYS.FLG.OnSignedOperation(res == 0, res>>(dstSize-1) == 1, overflow)
}

func (cpu *CPU) imul() {
	cpu._imathRS(func(a, b types.SValue) types.SValue { return a * b })
}

func (cpu *CPU) addf() {
	cpu._fmathRS(func(a, b float32) float32 { return a + b })
}

func (cpu *CPU) subf() {
	cpu._fmathRS(func(a, b float32) float32 { return a - b })
}

func (cpu *CPU) mulf() {
	cpu._fmathRS(func(a, b float32) float32 { return a * b })
}

func (cpu *CPU) divf() {
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcUVal := types.Value(cpu.castSrcToImm(source))
	dstUVal, err := cpu.RRAM.GetValue(dstRegID)

	dstVal := *((*float32)(unsafe.Pointer(&dstUVal)))
	srcVal := *((*float32)(unsafe.Pointer(&srcUVal)))
	if srcVal == 0 {
		cpu.SIGFPE()
		return
	}

	if err != nil {
		cpu.SIGILL()
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := dstVal / srcVal
	resUVal := *((*types.Value)(unsafe.Pointer(&res)))
	overflow := cpu.RRAM.PutValue(dstRegID, resUVal)
	cpu.RRAM.SYS.FLG.OnSignedOperation(res == 0, resUVal>>(dstSize-1) == 1, overflow)
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
		cpu.SIGILL()
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
	jump := cpu.getJump()
	*cpu.RRAM.SYS.NIR = cpu.castJumpToExeAddr(jump)
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) jif() {
	flag := cpu.getFlag()
	if cpu.RRAM.SYS.FLG.F(flag) {
		cpu.jmp()
	}
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) jnf() {
	flag := cpu.getFlag()
	if !cpu.RRAM.SYS.FLG.F(flag) {
		cpu.jmp()
	}
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) jiz() {
	regID := cpu.getReg()
	val, err := cpu.RRAM.GetValue(regID)
	if err != nil {
		cpu.SIGILL()
	}

	if val == 0 {
		cpu.jmp()
	}
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) jnz() {
	regID := cpu.getReg()
	val, err := cpu.RRAM.GetValue(regID)
	if err != nil {
		cpu.SIGILL()
	}

	if val != 0 {
		cpu.jmp()
	}
	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) lbl() {
	lblDst := cpu.getLabelDestination()
	isRegister, id := cpu.castLabelDstToModeAddr(lblDst)
	switch isRegister {
	case true:
		overflow := cpu.RRAM.PutValue(id, types.Value(*cpu.RRAM.SYS.NIR))
		if overflow {
			cpu.RRAM.SYS.FLG.FCOn()
		}
	case false:
		cpu.postLabelExeAddr(types.Address(id), *cpu.RRAM.SYS.NIR)
	}

	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) call() {
	*cpu.RRAM.SYS.MBR = *cpu.RRAM.SYS.NIR
	cpu.pushCallStack()

	jump := cpu.getJump()
	*cpu.RRAM.SYS.NIR = cpu.castJumpToExeAddr(jump)

	cpu.RRAM.SYS.FLG.Drop()
}

func (cpu *CPU) ret() {
	cpu.popCallStack()
	switch *cpu.RRAM.SYS.MBR == consts.MAX_WORD32 {
	case true:
		*cpu.RRAM.SYS.NIR = *cpu.RRAM.SYS.NIB
		cpu.RRAM.SYS.FLG = cpu.RRAM.SYS.FLB

		*cpu.RRAM.SYS.NIB = 0
		cpu.RRAM.SYS.FLB = 0

		cpu.OUTP.INTN = consts.SIGNONE
	case false:
		*cpu.RRAM.SYS.NIR = *cpu.RRAM.SYS.MBR
	}
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
		cpu.SIGILL()
		return
	}

	mask := types.Value((1 << dstSize) - 1)
	cpu.RRAM.PutValue(dstRegID, dstVal^mask)
}

func (cpu *CPU) int() {
	if !cpu.RRAM.SYS.FLG.FI() {
		return
	}

	cpu.OUTP.INT = true
	cpu.OUTP.INTN = byte(cpu.getInt())
}

func (cpu *CPU) _bitRR(fn func(a, b types.Value) types.Value) {
	dstRegID := cpu.getReg()
	srcRegID := cpu.getReg()

	dstVal, err1 := cpu.RRAM.GetValue(dstRegID)
	srcVal, err2 := cpu.RRAM.GetValue(srcRegID)
	if err1 != nil || err2 != nil {
		cpu.SIGILL()
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
		cpu.SIGILL()
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
		cpu.SIGILL()
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
		cpu.SIGILL()
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

func (cpu *CPU) _fmathRS(fn func(a, b float32) float32) {
	dstRegID := cpu.getReg()
	source := cpu.getSrc()

	srcUVal := types.Value(cpu.castSrcToImm(source))
	dstUVal, err := cpu.RRAM.GetValue(dstRegID)

	srcVal := *((*float32)(unsafe.Pointer(&srcUVal)))
	dstVal := *((*float32)(unsafe.Pointer(&dstUVal)))

	if err != nil {
		cpu.SIGILL()
		return
	}

	dstSize := cpu.RRAM.GetRegSize(dstRegID)
	res := fn(dstVal, srcVal)
	resUVal := *((*types.Value)(unsafe.Pointer(&res)))
	overflow := cpu.RRAM.PutValue(dstRegID, resUVal)
	cpu.RRAM.SYS.FLG.OnSignedOperation(res == 0, resUVal>>(dstSize-1) == 1, overflow)
}
