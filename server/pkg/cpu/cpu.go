package cpu

import (
	"errors"
	"fmt"
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
	_ = mem.NewSegment(consts.INT_SEG, consts.BiKB)
	_ = mem.NewSegment(consts.LBL_SEG, 3*consts.BiGB/8)
	_ = mem.NewSegment(consts.CLL_SEG, 3*consts.BiGB/8)
	cpu := CPU{
		RRAM: register.InitRRAM(),
		XMEM: &mem,
	}
	cpu.InitInterrupts()
	return cpu
}

const DefaultHandlersOffset = 21

func (cpu *CPU) InitInterrupts() {
	cpu.RRAM.SYS.FLG.FIOn()
	defaultHandlers := cpu.XMEM.At(consts.EXE_SEG).GetMaxAddr() - DefaultHandlersOffset

	// Default ignore exception handler
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0x0, consts.C_SKIP)
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0x4, consts.C_RET)

	// Default crush/wait exception handler
	cpu.XMEM.At(consts.EXE_SEG).SetWord32(defaultHandlers+0x8, consts.C_HLT)
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

	//maxIntAddr := cpu.XMEM.At(consts.INT_SEG).GetMaxAddr()
	//for i := 4*types.Address(consts.SIGILL) + 4; i < maxIntAddr; i += 4 {
	//	cpu.XMEM.At(consts.INT_SEG).SetWord32(i, types.Word32(defaultHandlers+0x8))
	//}
}

func (cpu *CPU) fetch(segmentID types.SegmentID, addr types.Address) {
	var err error
	*cpu.RRAM.SYS.MBR, err = cpu.XMEM.At(segmentID).GetWord32(addr)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
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

func (cpu *CPU) popCallStack() {
	*cpu.RRAM.SYS.CSP -= 4
	cpu.fetch(consts.CLL_SEG, types.Address(*cpu.RRAM.SYS.CSP))
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

func (cpu *CPU) pushCallStack() {
	*cpu.RRAM.SYS.MBR = *cpu.RRAM.SYS.NIR
	cpu.post(consts.CLL_SEG, types.Address(*cpu.RRAM.SYS.CSP), 32)
	*cpu.RRAM.SYS.CSP += 4
}

func (cpu *CPU) pushIntCallStack() {
	*cpu.RRAM.SYS.MBR = consts.MAX_WORD32
	cpu.post(consts.CLL_SEG, types.Address(*cpu.RRAM.SYS.CSP), 32)
	*cpu.RRAM.SYS.CSP += 4
}

func (cpu *CPU) DeclareLabel(label types.Address, exeAddr types.Word32) error {
	return cpu.XMEM.At(consts.LBL_SEG).SetWord32(label*4, exeAddr)
}

func (cpu *CPU) DeclareILabel(ilabel types.Address, exeAddr types.Word32) error {
	return cpu.XMEM.At(consts.INT_SEG).SetWord32(ilabel*4, exeAddr)
}

func (cpu *CPU) SetRegister(id types.Word32, value types.Value) {
	if !cpu.RRAM.PutValue(id, value) {
		cpu.SIGSEGV()
	}
}

func (cpu *CPU) SetMemory8(sid types.SegmentID, addr types.Address, value types.Word8) {
	if cpu.XMEM.At(sid).SetByte(addr, value) != nil {
		cpu.SIGSEGV()
	}
}

func (cpu *CPU) SetMemory16(sid types.SegmentID, addr types.Address, value types.Word16) {
	if cpu.XMEM.At(sid).SetWord16(addr, value) != nil {
		cpu.SIGSEGV()
	}
}

func (cpu *CPU) SetMemory32(sid types.SegmentID, addr types.Address, value types.Word32) {
	if cpu.XMEM.At(sid).SetWord32(addr, value) != nil {
		cpu.SIGSEGV()
	}
}

func (cpu *CPU) Tick() bool {
	*cpu.RRAM.SYS.IR = *cpu.RRAM.SYS.NIR
	*cpu.RRAM.SYS.NIR += 4

	halted := cpu.Exec()

	if cpu.OUTP.TERM {
		cpu.ForceSIGINT()
		return true
	}

	if cpu.InterruptCheck() {
		return halted
	}

	if cpu.RRAM.SYS.FLG.FT() && !halted && !cpu.OUTP.INTA {
		cpu.SIGTRACE()
		cpu.InterruptCheck()
		return halted
	}

	if cpu.OUTP.INTA && cpu.OUTP.INTN == consts.SIGNONE {
		cpu.OUTP.INTA = false
	}

	return halted
}

func (cpu *CPU) Exec() bool {
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
	case consts.C_RMD:
		cpu.rmd()
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
	case consts.C_EI:
		cpu.ei()
	case consts.C_DI:
		cpu.di()
	case consts.C_INT:
		cpu.int()
	case consts.C_HLT:
		cpu.skip()
		return true
	default:
		cpu.SIGILL()
	}
	return false
}
