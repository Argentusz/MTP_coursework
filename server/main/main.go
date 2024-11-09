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
	mtp.InitInterrupts()

	var program = []string{
		"ei",
		"mov rh0 0xFFFF",
		"mov rl0 0xFFF0",
		"div rx0 0",
		"add rb0 1",
		"mov [rx0] 1",
		"add rb0 1",
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
		finished = mtp.Tick()
		if StepByStep {
			mtp.MarshallHuman()
			fmt.Scanln(&str)
		}
	}

	mtp.MarshallHuman()
}
