package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/supervisor"
	"os"
	"strings"
)

func main() {
	base := flag.String("base", "../projects", "Base folder form which files will be opened to compile")
	intr := flag.Bool("intr", false, "Turns on/off interrupts by default")
	trace := flag.Bool("trace", false, "Turns on/off step-by-step mode")
	sudo := flag.Bool("sudo", false, "Turns on/off sudo mode")
	marshall := flag.String("marshall", "human", "Output state in human or JSON")
	flag.Parse()

	var marshallMode byte
	switch *marshall {
	case "human":
		marshallMode = supervisor.MarshallHuman
	case "JSON":
		marshallMode = supervisor.MarshallJSON
	default:
		panic(fmt.Sprintf("unknown marshall mode \"%s\"", *marshall))
	}

	visor := supervisor.InitSupervisor(*base, *intr, *trace, *sudo, marshallMode)

	onChange := func() {
		str, err := visor.Marshall()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
		fmt.Println(str)
	}

	stdin := bufio.NewReader(os.Stdin)
	for finished := false; !finished; {
		input, err := stdin.ReadString('\n')
		if err != nil {
			os.Stderr.WriteString(err.Error())
			continue
		}
		input = strings.Trim(input, "\n\t\r ")

		switch {
		case input == "quit":
			finished = true
		case strings.HasPrefix(input, "compile "):
			fileName := strings.TrimPrefix(input, "compile ")
			err := visor.Compile(fileName)
			onChange()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintln("[Error: Compilation]", err.Error()))
			}

		case input == "run" || input == "":
			go func() {
				err := visor.Run()
				onChange()
				if err != nil {
					os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
				}
			}()
		case input == "run all":
			go func() {
				err := visor.RunAll()
				onChange()
				if err != nil {
					os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
				}
			}()
		case input == "reset":
			visor.Reset()
			onChange()
		case input == "term":
			err := visor.Terminate()
			onChange()
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
			}
		case strings.HasPrefix(input, "trace "):
			onOff := strings.TrimPrefix(input, "trace ")
			switch onOff {
			case "on":
				visor.TraceOn()
				onChange()
			case "false":
				visor.TraceOff()
				onChange()
			default:
				onChange()
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] trace can be only on or off"))
			}

		case strings.HasPrefix(input, "sudo "):
			onOff := strings.TrimPrefix(input, "sudo ")
			switch onOff {
			case "on":
				visor.SudoOn()
				onChange()
			case "false":
				visor.SudoOff()
				onChange()
			default:
				onChange()
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] sudo can be only on or off"))
			}

		case strings.HasPrefix(input, "intr "):
			onOff := strings.TrimPrefix(input, "intr ")
			switch onOff {
			case "on":
				visor.IntrOn()
				onChange()
			case "false":
				visor.IntrOff()
				onChange()
			default:
				onChange()
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] intr can be only on or off"))
			}

		default:
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] Command \"%s\" not found\n", input))
		}
	}

}

//func main() {
//	const StepByStep = false
//	mtp := cpu.InitCPU()
//
//	fileContent, err := os.ReadFile("../projects/examples/primes.mtp")
//	if err != nil {
//		panic(err.Error())
//	}
//
//	compiled, err := compiler.Compile(strings.Split(string(fileContent), "\n"))
//
//	if err != nil {
//		panic(err.Error())
//	}
//
//	program := compiled.Output
//	labels := compiled.Labels
//	ilabels := compiled.ILabels
//
//	marshallExeSeg := func() {
//		for i, v := range compiled.Input {
//			if int(*mtp.RRAM.SYS.NIB) == i*4 {
//				fmt.Print("[*] ")
//			} else {
//				fmt.Print("    ")
//			}
//
//			f1 := false
//			for label, addr := range compiled.Labels {
//				if int(addr) == i*4 {
//					fmt.Printf("<0x%04x>: ", label)
//					f1 = true
//					break
//				}
//			}
//
//			f2 := false
//			for ilabel, addr := range compiled.ILabels {
//				if int(addr) == i*4 {
//					fmt.Printf("{0x%04x}: ", ilabel)
//					f2 = true
//					break
//				}
//			}
//
//			if !f1 && !f2 {
//				fmt.Print("          ")
//			}
//
//			fmt.Println(v)
//		}
//	}
//
//	for i, line := range program {
//		err = mtp.XMEM.At(consts.EXE_SEG).SetWord32(types.Address(i*4), line)
//		if err != nil {
//			panic(err.Error())
//		}
//	}
//
//	for label, addr := range labels {
//		err = mtp.DeclareLabel(types.Address(label), addr)
//		if err != nil {
//			panic(err.Error())
//		}
//	}
//
//	for ilabel, addr := range ilabels {
//		err = mtp.DeclareILabel(types.Address(ilabel), addr)
//		if err != nil {
//			panic(err.Error())
//		}
//	}
//
//	if StepByStep {
//		mtp.RRAM.SYS.FLG.FTOn()
//	}
//
//	mtp.MarshallHuman()
//	marshallExeSeg()
//
//	var str string
//	fmt.Scanln(&str)
//
//	for finish := str == "q"; !finish; {
//		halted := mtp.Tick()
//		if !halted {
//			continue
//		}
//
//		finish = true
//		mtp.MarshallHuman()
//		marshallExeSeg()
//		if mtp.OUTP.INTA && mtp.OUTP.INTN == consts.SIGTRACE {
//			//time.Sleep(250 * time.Millisecond)
//			//finish = false
//			fmt.Scanln(&str)
//			finish = str == "q"
//		}
//	}
//
//	switch mtp.OUTP.INTN {
//	case consts.SIGNONE:
//		fmt.Println("Program finished successfully")
//	case consts.SIGFPE:
//		fmt.Println("Program interrupted: SIGFPE")
//	case consts.SIGTRACE:
//		fmt.Println("Program interrupted: SIGTRACE")
//	case consts.SIGSEGV:
//		fmt.Println("Program interrupted: SIGSEGV")
//	case consts.SIGTERM:
//		fmt.Println("Program interrupted: SIGTERM")
//	case consts.SIGINT:
//		fmt.Println("Program interrupted: SIGINT")
//	case consts.SIGIIE:
//		fmt.Println("Program interrupted: SIGIIE")
//	case consts.SIGILL:
//		fmt.Println("Program interrupted: SIGILL")
//	default:
//		fmt.Println("Program interrupted with error code", mtp.OUTP.INTN)
//	}
//
//}
