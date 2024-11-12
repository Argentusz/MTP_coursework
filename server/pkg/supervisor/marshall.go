package supervisor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

const (
	MarshallHuman byte = iota
	MarshallJSON
)

type MarshalledSystem struct {
	Running bool
	Input   []string
	CPU     cpu.CPU
}

func (s *Supervisor) Marshall() (string, error) {
	switch s.marshallMode {
	case MarshallHuman:
		return s.marshallHuman()
	case MarshallJSON:
		return s.marshallJSON()
	}

	return "", errors.New(fmt.Sprintf("marshall mode %d not found", s.marshallMode))
}

func (s *Supervisor) marshallJSON() (string, error) {
	var ms MarshalledSystem
	ms.Running = s.running
	ms.CPU = s.cpu
	ms.Input = s.compiler.Input

	bytes, err := json.Marshal(ms)
	return string(bytes), err
}

func (s *Supervisor) marshallHuman() (string, error) {
	str := ""

	str += fmt.Sprintln("Memory:")
	xmem, err := s.marshallHumanXMEM()
	if err != nil {
		return "", err
	}
	str += xmem

	str += fmt.Sprintln("Registers:")
	str += s.marshallHumanRRAM()

	str += fmt.Sprintln("Wires:")
	str += s.marshallHumanOUTP()

	str += fmt.Sprintln("Program:")
	str += s.marshallHumanProgram()

	return str, nil
}

func (s *Supervisor) marshallHumanXMEM() (string, error) {
	str := fmt.Sprintf("| %7s | %33s | %35s |\n", "ADDR", "EXE_SEG", "USR_SEG")
	for i := types.Address(0); i < 10*4; i += 4 {
		exew, erre := s.cpu.XMEM.At(consts.EXE_SEG).GetWord32(i)
		usrw, erru := s.cpu.XMEM.At(consts.USR_SEG).GetWord32(i)
		if erre != nil {
			return "", erre
		}
		if erru != nil {
			return "", erru
		}

		str += fmt.Sprintf("| 0x%05x | %026b %06b | %08b %08b %08b %08b |\n", i, exew>>6, exew&0b111111, (usrw>>24)&consts.MAX_WORD8, (usrw>>16)&consts.MAX_WORD8, (usrw>>8)&consts.MAX_WORD8, usrw&consts.MAX_WORD8)
	}
	return str, nil
}

func (s *Supervisor) marshallHumanRRAM() string {
	str := ""
	str += fmt.Sprintf("| %17s |\n", "SYS")
	str += fmt.Sprintf("---------------------\n")
	str += fmt.Sprintf("| %4s | 0x%08x |\n", "IR", *s.cpu.RRAM.SYS.IR)
	str += fmt.Sprintf("| %4s | %10d |\n", "NIB", *s.cpu.RRAM.SYS.NIB)
	str += fmt.Sprintf("| %4s | %10d |\n", "MBR", *s.cpu.RRAM.SYS.MBR)
	str += fmt.Sprintf("| %4s | %10d |\n", "TMP", *s.cpu.RRAM.SYS.TMP)
	str += fmt.Sprintf("| %4s |   %08b |\n", "FLG", s.cpu.RRAM.SYS.FLG)
	str += fmt.Sprintf("| %4s |   %08b |\n\n", "FLB", s.cpu.RRAM.SYS.FLB)

	str += fmt.Sprintf("| %10s |      | %60s%17s |     | %17s |\n", "SGPR", " ", "GPR", "XGPR")
	str += fmt.Sprintf("%14s      %81s     %21s\n", "--------------", "---------------------------------------------------------------------------------", "---------------------")

	for i := 0; i < 8; i++ {
		str += fmt.Sprintf("| %4s | %3d |      ", fmt.Sprintf("rb%d", i), *s.cpu.RRAM.SGPRs[i])
		str += fmt.Sprintf("| %4s | %10d |", fmt.Sprintf("rw%d", i), *s.cpu.RRAM.GPRs[i])
		str += fmt.Sprintf(" %4s | %10d |", fmt.Sprintf("rw%d", i+8), *s.cpu.RRAM.GPRs[i+8])
		str += fmt.Sprintf(" %4s | %10d |", fmt.Sprintf("rw%d", i+16), *s.cpu.RRAM.GPRs[i+16])
		str += fmt.Sprintf(" %4s | %10d |     ", fmt.Sprintf("rw%d", i+24), *s.cpu.RRAM.GPRs[i+24])
		str += fmt.Sprintf("| %4s | %10d |\n", fmt.Sprintf("rx%d", i), *s.cpu.RRAM.XGPRs[i])
	}

	return str
}

func (s *Supervisor) marshallHumanOUTP() string {
	return fmt.Sprintf("INT = %v, INTA = %v, INTN = %v\n",
		s.cpu.OUTP.INT, s.cpu.OUTP.INTA, s.cpu.OUTP.INTN)
}

func (s *Supervisor) marshallHumanProgram() string {
	str := ""
	for i, v := range s.compiler.Input {
		if int(*s.cpu.RRAM.SYS.NIB) == i*4 {
			str += fmt.Sprint("[*] ")
		} else {
			str += fmt.Sprint("    ")
		}

		f1 := false
		for label, addr := range s.compiler.Labels {
			if int(addr) == i*4 {
				str += fmt.Sprintf("<0x%04x>: ", label)
				f1 = true
				break
			}
		}

		f2 := false
		for ilabel, addr := range s.compiler.ILabels {
			if int(addr) == i*4 {
				str += fmt.Sprintf("{0x%04x}: ", ilabel)
				f2 = true
				break
			}
		}

		if !f1 && !f2 {
			str += fmt.Sprint("          ")
		}

		str += fmt.Sprintln(v)
	}

	return str
}
