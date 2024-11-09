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
	OUTP Outputs
}

func InitCPU() CPU {
	mem := xmem.InitExternalMemory()
	_ = mem.NewSegment(consts.EXE_SEG, 1.5*consts.BiGB)
	_ = mem.NewSegment(consts.USR_SEG, 2.0*consts.BiMB)
	_ = mem.NewSegment(consts.INT_SEG, 9*consts.BiKB)
	_ = mem.NewSegment(consts.LBL_SEG, 3*consts.BiGB/8)
	return CPU{
		RRAM: register.InitRRAM(),
		XMEM: &mem,
	}
}

func (cpu *CPU) InitInterrupts() {
	maxExeAddr := cpu.XMEM.At(consts.EXE_SEG).GetMaxAddr()
	defaultHandlers := maxExeAddr - 21

	// Default ignore exception handler
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0x0, consts.C_SKIP)
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0x4, consts.C_RET)

	// Default crush/wait exception handler
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0x8, consts.C_HALT)
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0xc, consts.C_RET)

	// Fill default Interrupt Descriptor Table
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGNONE), types.Word32(defaultHandlers))
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGFPE), types.Word32(defaultHandlers))

	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGTRACE), types.Word32(defaultHandlers+0x8))
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGSEGV), types.Word32(defaultHandlers+0x8))
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGTERM), types.Word32(defaultHandlers+0x8))
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGINT), types.Word32(defaultHandlers+0x8))
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGIIE), types.Word32(defaultHandlers+0x8))
	cpu.XMEM.At(consts.INT_SEG).SetWord32(4*types.Address(consts.SIGILL), types.Word32(defaultHandlers+0x8))

	maxIntAddr := cpu.XMEM.At(consts.INT_SEG).GetMaxAddr()
	for i := 4*types.Address(consts.SIGILL) + 4; i < maxIntAddr; i += 4 {
		cpu.XMEM.At(consts.INT_SEG).SetWord32(i, types.Word32(defaultHandlers+0x8))
	}
}

func (cpu *CPU) fetch(segmentID types.SegmentID, addr types.Address) {
	var err error
	*cpu.RRAM.SYS.MBR, err = cpu.XMEM.At(segmentID).GetWord32(addr)
	if err != nil {
		cpu.SIGSEGV()
	}
}

func (cpu *CPU) fetchInstr() {
	cpu.fetch(consts.EXE_SEG, types.Address(*cpu.RRAM.SYS.IR))
	*cpu.RRAM.SYS.TMP = *cpu.RRAM.SYS.MBR
}

func (cpu *CPU) fetchUsrData(addr types.Address) {
	cpu.fetch(consts.USR_SEG, addr)
}

func (cpu *CPU) fetchLabelExeAddr(label types.Address) {
	cpu.fetch(consts.LBL_SEG, label*4)
}

func (cpu *CPU) fetchIntExeAddr(intn types.Address) {
	cpu.fetch(consts.INT_SEG, intn*4)
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
		cpu.SIGSEGV()
	}
}

func (cpu *CPU) postUsrData(addr types.Address, val types.Word32, size byte) {
	*cpu.RRAM.SYS.MBR = val
	cpu.post(consts.USR_SEG, addr, size)
}

func (cpu *CPU) postLabelExeAddr(label types.Address, exeAddr types.Word32) {
	*cpu.RRAM.SYS.MBR = exeAddr
	cpu.post(consts.LBL_SEG, label*4, 32)
}

func (cpu *CPU) Tick() bool {
	*cpu.RRAM.SYS.NIR += 4
	cpu.fetchInstr()
	operator := cpu.getOperator()
	switch operator {
	case consts.C_SKIP:
		cpu.skip()
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
	case consts.C_DIV:
		cpu.div()
	case consts.C_IMOV:
		cpu.imov()
	case consts.C_IADD:
		cpu.iadd()
	case consts.C_IADC:
		cpu.iadc()
	case consts.C_ISUB:
		cpu.isub()
	case consts.C_ISBB:
		cpu.isbb()
	case consts.C_IMUL:
		cpu.imul()
	case consts.C_IDIV:
		cpu.idiv()
	case consts.C_ADDF:
		cpu.addf()
	case consts.C_SUBF:
		cpu.subf()
	case consts.C_MULF:
		cpu.mulf()
	case consts.C_DIVF:
		cpu.divf()
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
		cpu.jmp()
	case consts.C_JIF:
		cpu.jif()
	case consts.C_JNF:
		cpu.jnf()
	case consts.C_JIZ:
		cpu.jiz()
	case consts.C_JNZ:
		cpu.jnz()
	case consts.C_LBL:
		cpu.lbl()
	case consts.C_CALL:
		cpu.call()
	case consts.C_RET:
		cpu.ret()
	case consts.C_HALT:
		cpu.skip()
	case consts.C_EI:
		cpu.ei()
	case consts.C_DI:
		cpu.di()
	case consts.C_INT:
		cpu.int()
	default:
		cpu.SIGILL()
		panic("unknown operator")
		return false
	}

	cpu.InterruptCheck()
	*cpu.RRAM.SYS.IR = *cpu.RRAM.SYS.NIR
	return operator == consts.C_HALT
}
