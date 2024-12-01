package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/supervisor"
	"github.com/Argentusz/MTP_coursework/pkg/types"
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
	onChange(&visor)
	mainLoop(&visor)
}

func mainLoop(visor *supervisor.Supervisor) {
	stdin := bufio.NewReader(os.Stdin)
	for finished := false; !finished; {
		input, err := stdin.ReadString('\n')
		if err != nil {
			os.Stderr.WriteString(err.Error())
			continue
		}

		input = strings.Trim(input, "\n\t\r ")
		finished = do(input, visor)
	}
}

func onChange(visor *supervisor.Supervisor) {
	str, err := visor.Marshall()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(str + "\x00")
}

func do(input string, visor *supervisor.Supervisor) bool {
	switch {

	case input == "quit":
		return true

	case strings.HasPrefix(input, "compile "):
		fileName := strings.TrimPrefix(input, "compile ")
		err := visor.Compile(fileName)
		onChange(visor)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintln("[Error: Compilation]", err.Error()))
		}

	case input == "run" || input == "":
		go func() {
			err := visor.Run()
			onChange(visor)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
			}
		}()

	case input == "run all":
		go func() {
			err := visor.RunAll()
			onChange(visor)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
			}
		}()

	case input == "reset":
		visor.Reset()
		onChange(visor)

	case input == "term":
		err := visor.Terminate()
		onChange(visor)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
		}

	case strings.HasPrefix(input, "set "):
		split := strings.SplitN(input, " ", 2)
		if len(split) != 2 {
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", "set expected JSON"))
			break
		}
		jsonStr := split[1]
		updates := struct {
			RRAM map[types.Word32]types.Value
			XMEM map[types.SegmentID]map[types.Address]types.Word8
		}{}
		err := json.Unmarshal([]byte(jsonStr), &updates)
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
			break
		}
		for rid, val := range updates.RRAM {
			err := visor.SetRegister(rid, val)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
				return false
			}
		}
		for sid, m := range updates.XMEM {
			for addr, val := range m {
				err := visor.SetMemory8(sid, addr, val)
				if err != nil {
					os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor]", err.Error()))
					return false
				}
			}
		}
		onChange(visor)

	case strings.HasPrefix(input, "trace "):
		onOff := strings.TrimPrefix(input, "trace ")
		switch onOff {
		case "on":
			visor.TraceOn()
			onChange(visor)
		case "off":
			visor.TraceOff()
			onChange(visor)
		default:
			onChange(visor)
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] trace can be only on or off"))
		}

	case strings.HasPrefix(input, "sudo "):
		onOff := strings.TrimPrefix(input, "sudo ")
		switch onOff {
		case "on":
			visor.SudoOn()
			onChange(visor)
		case "off":
			visor.SudoOff()
			onChange(visor)
		default:
			onChange(visor)
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] sudo can be only on or off"))
		}

	case strings.HasPrefix(input, "intr "):
		onOff := strings.TrimPrefix(input, "intr ")
		switch onOff {
		case "on":
			visor.IntrOn()
			onChange(visor)
		case "off":
			visor.IntrOff()
			onChange(visor)
		default:
			onChange(visor)
			os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] intr can be only on or off"))
		}

	default:
		onChange(visor)
		os.Stderr.WriteString(fmt.Sprintln("[Error: Supervisor] Command \"%s\" not found\n", input))
	}
	return false
}
