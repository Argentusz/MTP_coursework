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
	compiled, err := compiler.Compile([]string{
		"; Программа для нахождения всех делителей числа",
		"#start $main",
		"$define $input 256",
		"$define $write 1",
		"$define $just_test 2",
		"$define $for 3",
		"$define $repeat jnz rx2 $for",
		"",
		"",
		"@label $write ; Запись делителя rx1 числа rx0",
		"    mov [rw0] rx1",
		"    add rw0 4",
		"    ret",
		"",
		"",
		"@label $main",
		"    #mov32 rx0 $input   ; Исходное число",
		"    mov rw0 0           ; Адрес результата",
		"    ",
		"    lbl $for",
		"    int $sigtrace",
		"    add rx1 1",
		"    mov rx2 rx0",
		"    rmd rx2 rx1",
		"    $repeat",
		"    call $write",
		"    mov rx2 rx0",
		"    mov rx3 rx1",
		"    sub rx2 rx3",
		"    $repeat",
		"    ",
		"    #mov32 rx1 $m32",
		"    mov [rw0] rx1",
		"    ",
		"    call $just_test",
		"    ",
		"    hlt",
		"",
		"",
		"@label $just_test",
		"    #mov32 rw0 f1.5",
		"    #mov32 rw1 f2.3",
		"    addf rw0 rw1",
		"    ret",
		"",
		"",
		"@ilabel $sigtrace",
		"    add rw10 1",
		"    ret",
		"",
	})

	if err != nil {
		panic(err.Error())
	}

	program := compiled.Output
	labels := compiled.Labels
	ilabels := compiled.ILabels

	for i, line := range program {
		err = mtp.XMEM.At(consts.EXE_SEG).SetWord32(types.Address(i*4), line)
		if err != nil {
			panic(err.Error())
		}
	}

	for label, addr := range labels {
		err = mtp.DeclareLabel(types.Address(label), addr)
		if err != nil {
			panic(err.Error())
		}
	}

	for ilabel, addr := range ilabels {
		err = mtp.DeclareILabel(types.Address(ilabel), addr)
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
