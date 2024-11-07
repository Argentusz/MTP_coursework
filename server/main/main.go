package main

import (
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/interpreter"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func main() {
	const StepByStep = false
	mtp := cpu.InitCPU()

	var program = []string{
		"mov rh0 0x3fc0",
		"mov rl0 0x0000", // rx0 = 1.5
		"skip",
		"mov rh1 0x404c",
		"mov rl1 0xcccd", // rx1 ≈ 3.2
		"skip",
		"subf rx1 rx0",
		//"mov rw1 16",
		//"mov rw2 1",
		//"mov rx1 12",
		//"lbl rx0",
		//"mul rw2 2",
		//"sub rw1 1",
		//"jnz rw1 [rx0]",
		//"mov [rx1] rw2",
		"halt",
	}

	for i, line := range program {
		compiled, err := interpreter.Convert(line)
		if err != nil {
			panic(err.Error())
		}

		err = mtp.XMEM.At(consts.EXE_SEG).SetWord32(types.Address(i*4), compiled)
		if err != nil {
			panic(err.Error())
		}
	}

	str := ""
	finished := false
	for !finished && str != ":q" {
		finished = mtp.Exec()
		if StepByStep {
			mtp.MarshallHuman()
			fmt.Scanln(&str)
		}
	}

	mtp.MarshallHuman()
}
