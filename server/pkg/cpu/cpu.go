package cpu

import (
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

func (cpu *CPU) post(segmentID types.SegmentID, addr types.Address) {
	err := cpu.XMEM.At(segmentID).SetWord32(addr, *cpu.RRAM.SYS.MBR)
	if err != nil {
		// SIGSEGV
	}
}

func (cpu *CPU) postUsrData(addr types.Address, val types.Word32) {
	*cpu.RRAM.SYS.MBR = val
	cpu.post(consts.USR_SEG, addr)
}

func (cpu *CPU) Exec() bool {
	cpu.fetchInstr()
	switch cpu.getOperator() {
	case consts.C_ADD:
		cpu.add()
	case consts.C_MOV:
		cpu.mov()
	default:
		// SIGILL
		return false
	}

	*cpu.RRAM.SYS.IR += 4
	return true
}
