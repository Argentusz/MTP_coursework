package cpu

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) skip() {

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

	cpu.postUsrData(types.Address(dstAddr), types.Word32(srcImm))
}

func (cpu *CPU) add() {
	param1 := cpu.getReg()
	param2 := cpu.getSrc()

	src := types.Value(cpu.castSrcToImm(param2))
	dest, err := cpu.RRAM.GetValue(param1)
	if err != nil {
		// SIGILL
		return
	}

	cpu.RRAM.PutValue(param1, src+dest)
}
