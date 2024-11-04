package main

import (
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/interpreter"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func main() {
	mtp := cpu.InitCPU()

	var program = []string{
		"mov rw1 1",
		"mov rw2 2",
		"mov rw3 0",
		"add rw1 rw2",
		"mov [rw3] rw1",
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

	for i := 0; i < len(program); i++ {
		fmt.Println(program[i])
		mtp.Exec()
	}

	mtp.MarshallHuman()
}
