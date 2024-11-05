package cpu

import (
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func (cpu *CPU) MarshallHuman() {
	fmt.Println("Memory:")
	cpu.marshallHumanXMEM()

	fmt.Println("Registers:")
	cpu.marshallHumanRRAM()
}

func (cpu *CPU) marshallHumanXMEM() {
	fmt.Printf("| %7s | %33s | %35s |\n", "ADDR", "EXE_SEG", "USR_SEG")
	for i := types.Address(0); i < 5*4; i += 4 {
		exew, erre := cpu.XMEM.At(consts.EXE_SEG).GetWord32(i)
		usrw, erru := cpu.XMEM.At(consts.USR_SEG).GetWord32(i)
		if erre != nil {
			panic(erre.Error())
		}
		if erru != nil {
			panic(erru.Error())
		}

		fmt.Printf("| 0x%05x | %027b %05b | %08b %08b %08b %08b |\n", i, exew>>5, exew&0b11111, (usrw>>24)&consts.MAX_WORD8, (usrw>>16)&consts.MAX_WORD8, (usrw>>8)&consts.MAX_WORD8, usrw&consts.MAX_WORD8)
	}
}

// |      SYS |
// | IR   | 0 |
// | IRB  | 0 |
// | MBR  | 0 |
// | TMP  | 0 |
// | FLG  | 0 |
// | FLGI | 0 |
func (cpu *CPU) marshallHumanRRAM() {

	fmt.Printf("| %17s |\n", "SYS")
	fmt.Printf("---------------------\n")
	fmt.Printf("| %4s | %10d |\n", "IR", *cpu.RRAM.SYS.IR)
	fmt.Printf("| %4s | %10d |\n", "IRB", *cpu.RRAM.SYS.IRB)
	fmt.Printf("| %4s | %10d |\n", "MBR", *cpu.RRAM.SYS.MBR)
	fmt.Printf("| %4s | %10d |\n", "TMP", *cpu.RRAM.SYS.TMP)
	fmt.Printf("| %4s |   %08b |\n", "FLG", cpu.RRAM.SYS.FLG)
	fmt.Printf("| %4s |   %08b |\n\n", "FLGI", cpu.RRAM.SYS.FLGI)

	fmt.Printf("| %10s |      | %60s%17s |     | %17s |\n", "SGPR", " ", "GPR", "XGPR")
	fmt.Printf("%14s      %81s     %21s\n", "--------------", "---------------------------------------------------------------------------------", "---------------------")

	for i := 0; i < 8; i++ {
		fmt.Printf("| %4s | %3d |      ", fmt.Sprintf("rb%d", i), *cpu.RRAM.SGPRs[i])
		fmt.Printf("| %4s | %10d |", fmt.Sprintf("rw%d", i), *cpu.RRAM.GPRs[i])
		fmt.Printf(" %4s | %10d |", fmt.Sprintf("rw%d", i+8), *cpu.RRAM.GPRs[i+8])
		fmt.Printf(" %4s | %10d |", fmt.Sprintf("rw%d", i+16), *cpu.RRAM.GPRs[i+16])
		fmt.Printf(" %4s | %10d |     ", fmt.Sprintf("rw%d", i+24), *cpu.RRAM.GPRs[i+24])
		fmt.Printf("| %4s | %10d |\n", fmt.Sprintf("rx%d", i), *cpu.RRAM.XGPRs[i])
	}

}
