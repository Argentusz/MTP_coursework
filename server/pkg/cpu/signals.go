package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) InterruptCheck() {
	if !cpu.OUTP.INT {
		return
	}

	cpu.fetchIntExeAddr(types.Address(cpu.OUTP.INTN))

	cpu.OUTP.INT = false
	cpu.OUTP.INTA = true
	*cpu.RRAM.SYS.NIB = *cpu.RRAM.SYS.NIR
	*cpu.RRAM.SYS.NIR = *cpu.RRAM.SYS.MBR
}

func (cpu *CPU) setInterrupt(intn byte) {
	if !cpu.RRAM.SYS.FLG.FI() || cpu.OUTP.INTN != consts.SIGNONE {
		return
	}

	cpu.OUTP.INT = true
	cpu.OUTP.INTN = intn
}

func (cpu *CPU) SIGFPE() {
	cpu.setInterrupt(consts.SIGFPE)
}

func (cpu *CPU) SIGTRACE() {
	cpu.setInterrupt(consts.SIGTRACE)
}

func (cpu *CPU) SIGSEGV() {
	cpu.setInterrupt(consts.SIGSEGV)
}

func (cpu *CPU) SIGTERM() {
	cpu.setInterrupt(consts.SIGTERM)
}

func (cpu *CPU) SIGINT() {
	cpu.setInterrupt(consts.SIGINT)
}

func (cpu *CPU) SIGIIE() {
	cpu.setInterrupt(consts.SIGIIE)
}

func (cpu *CPU) SIGILL() {
	cpu.setInterrupt(consts.SIGILL)
}
