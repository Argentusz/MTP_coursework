package main

import (
	"encoding/json"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/interpreter"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

func main() {
	mtp := cpu.InitCPU()
	var jsonOut = func() {
		bytes, err := json.Marshal(mtp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(bytes))
	}

	var memoryOut = func(seg types.SegmentID, from, to types.Address) {
		for i := from; i < to; i += 4 {
			v, _ := mtp.XMEM.At(seg).GetWord32(i)
			fmt.Printf("%032b\n", v)
		}
	}

	jsonOut()

	fmt.Println(mtp.XMEM.At(consts.EXE_SEG).GetByte(0))
	fmt.Println(mtp.XMEM.At(consts.EXE_SEG).SetByte(0, 1))
	fmt.Println(mtp.XMEM.At(consts.EXE_SEG).GetByte(0))

	jsonOut()

	var program = []string{
		"mov rw1 1",
		"mov rw2 2",
		"mov rw3 0",
		"add rw1 rw2",
		"mov [rw3] rw1",
	}

	for i, line := range program {
		compiled, err := interpreter.Convert(line)
		fmt.Println(compiled, err)
		if err != nil {
			panic(err.Error())
		}

		err = mtp.XMEM.At(consts.EXE_SEG).SetWord32(types.Address(i*4), compiled)
		if err != nil {
			panic(err.Error())
		}
	}

	memoryOut(consts.EXE_SEG, 0, 20)

	flag := true
	for flag {
		flag = mtp.Exec()
		jsonOut()
	}

	memoryOut(consts.USR_SEG, 0, 20)

}
