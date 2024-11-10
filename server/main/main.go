package main

import (
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/compiler"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func main() {
	const StepByStep = false
	mtp := cpu.InitCPU()
	program, err := compiler.Compile([]string{
		"mov rb0 10 ; Счётчик",
		"mov rw1 1",
		"lbl 1",
		"; Тело цикла",
		"mul rw1 2",
		"sub rb0 1 ; Уменьшаем счётчик",
		"jnz rb0 1",
		"add rb1 1",
		"hlt",
	}, &mtp)

	if err != nil {
		panic(err.Error())
	}

	for i, line := range program {
		err = mtp.XMEM.At(consts.EXE_SEG).SetWord32(types.Address(i*4), line)
		if err != nil {
			panic(err.Error())
		}
	}

	if StepByStep {
		mtp.RRAM.SYS.FLG.FTOn()
	}

	for finish := false; !finish; {
		halted := mtp.Tick()
		if !halted {
			continue
		}

		finish = true
		mtp.MarshallHuman()
		if mtp.OUTP.INTA && mtp.OUTP.INTN == consts.SIGTRACE {
			var str string
			fmt.Scanln(&str)
			finish = str == "q"
		}
	}

	switch mtp.OUTP.INTN {
	case consts.SIGNONE:
		fmt.Println("Program finished successfully")
	case consts.SIGFPE:
		fmt.Println("Program interrupted: SIGFPE")
	case consts.SIGTRACE:
		fmt.Println("Program interrupted: SIGTRACE")
	case consts.SIGSEGV:
		fmt.Println("Program interrupted: SIGSEGV")
	case consts.SIGTERM:
		fmt.Println("Program interrupted: SIGTERM")
	case consts.SIGINT:
		fmt.Println("Program interrupted: SIGINT")
	case consts.SIGIIE:
		fmt.Println("Program interrupted: SIGIIE")
	case consts.SIGILL:
		fmt.Println("Program interrupted: SIGILL")
	default:
		fmt.Println("Program interrupted with error code", mtp.OUTP.INTN)
	}

}
