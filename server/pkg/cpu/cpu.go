package cpu

import (
	"errors"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/register"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"github.com/Argentusz/MTP_coursework/pkg/xmem"
)

type CPU struct {
	RRAM register.RRAM
	XMEM *xmem.ExternalMemory
}

func InitCPU() CPU {
	mem := xmem.InitExternalMemory()
	_ = mem.NewSegment(consts.EXE_SEG, 1.5*consts.BiGB)
	_ = mem.NewSegment(consts.INT_SEG, 0.5*consts.BiGB)
	_ = mem.NewSegment(consts.USR_SEG, 2.0*consts.BiGB)
	return CPU{
		RRAM: register.InitRRAM(),
		XMEM: &mem,
	}
}

func (cpu *CPU) fetch(segmentID types.SegmentID, addr types.Address) {
	var err error
	*cpu.RRAM.SYS.MBR, err = cpu.XMEM.At(segmentID).GetWord32(addr)
	if err != nil {
		// SIGSEGV
	}
}

func (cpu *CPU) fetchInstr() {
	cpu.fetch(consts.EXE_SEG, types.Address(*cpu.RRAM.SYS.IR))
	*cpu.RRAM.SYS.TMP = *cpu.RRAM.SYS.MBR
}

func (cpu *CPU) fetchUsrData(addr types.Address) {
	cpu.fetch(consts.USR_SEG, addr)
}

func (cpu *CPU) post(segmentID types.SegmentID, addr types.Address, size byte) {
	var err error
	switch size {
	case 8:
		err = cpu.XMEM.At(segmentID).SetByte(addr, types.Word8(*cpu.RRAM.SYS.MBR))
	case 16:
		err = cpu.XMEM.At(segmentID).SetWord16(addr, types.Word16(*cpu.RRAM.SYS.MBR))
	case 32:
		err = cpu.XMEM.At(segmentID).SetWord32(addr, *cpu.RRAM.SYS.MBR)
	default:
		err = errors.New("can write only 8, 16 or 32-bit messages")
	}
	if err != nil {
		// SIGSEGV
	}
}

func (cpu *CPU) postUsrData(addr types.Address, val types.Word32, size byte) {
	*cpu.RRAM.SYS.MBR = val
	cpu.post(consts.USR_SEG, addr, size)
}

func (cpu *CPU) Exec() bool {
	cpu.fetchInstr()
	//fmt.Printf("Executing %027b %05b\n", *cpu.RRAM.SYS.TMP>>5, *cpu.RRAM.SYS.TMP&0b11111)
	operator := cpu.getOperator()
	//fmt.Printf("Operator code: %05b\n", operator)
	//fmt.Printf("TMP state: %032b\n", *cpu.RRAM.SYS.TMP)
	switch operator {
	case consts.C_MOV:
		cpu.mov()
	case consts.C_ADD:
		cpu.add()
	case consts.C_ADC:
		cpu.adc()
	case consts.C_SUB:
		cpu.sub()
	case consts.C_SBB:
		cpu.sbb()
	case consts.C_MUL:
		cpu.mul()
	case consts.C_IMUL:
		cpu.imul()
	case consts.C_DIV:
		cpu.div()
	case consts.C_SHL:
		cpu.shl()
	case consts.C_SHR:
		cpu.shr()
	case consts.C_SAR:
		cpu.sar()
	case consts.C_AND:
		cpu.and()
	case consts.C_OR:
		cpu.or()
	case consts.C_XOR:
		cpu.xor()
	case consts.C_NOT:
		cpu.not()
	case consts.C_JMP:
		panic("jmp is NOI")
	case consts.C_CALL:
		panic("call is NOI")
	case consts.C_RET:
		panic("ret is NOI")
	case consts.C_HALT:
		panic("halt is NOI")
	case consts.C_EI:
		panic("ei is NOI")
	case consts.C_DI:
		panic("di is NOI")
	case consts.C_INT:
		panic("int is NOI")
	case consts.C_ADDF:
		panic("addf is NOI")
	case consts.C_SUBF:
		panic("subf is NOI")
	case consts.C_MULF:
		panic("mulf is NOI")
	case consts.C_DIVF:
		panic("divf is NOI")
	default:
		// SIGILL
		return false
	}

	*cpu.RRAM.SYS.IR += 4
	return true
}
