package supervisor

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/compiler"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"os"
	"strings"
)

type Supervisor struct {
	cpu          cpu.CPU
	compiler     compiler.Compiler
	baseDir      string
	running      bool
	ableToRun    bool
	marshallMode byte
	flags        struct {
		intr  bool
		trace bool
		sudo  bool
	}
}

func InitSupervisor(baseDir string, intr, trace, sudo bool, marshallMode byte) Supervisor {
	var s Supervisor
	s.cpu = cpu.InitCPU()
	s.baseDir = baseDir
	s.marshallMode = marshallMode
	s.flags.intr, s.flags.trace, s.flags.sudo = intr, trace, sudo
	s.setFlags()
	return s
}

func (s *Supervisor) Compile(fileName string) error {
	s.Reset()
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", s.baseDir, fileName))
	if err != nil {
		return err
	}

	compiled, err := compiler.Compile(strings.Split(string(fileContent), "\n"))
	if err != nil {
		return err
	}

	program := compiled.Output
	labels := compiled.Labels
	ilabels := compiled.ILabels

	for i, line := range program {
		err = s.cpu.XMEM.At(consts.EXE_SEG).SetWord32(types.Address(i*4), line)
		if err != nil {
			return err
		}
	}

	for label, addr := range labels {
		err = s.cpu.DeclareLabel(types.Address(label), addr)
		if err != nil {
			return err
		}
	}

	for ilabel, addr := range ilabels {
		err = s.cpu.DeclareILabel(types.Address(ilabel), addr)
		if err != nil {
			return err
		}
	}

	s.compiler = compiled
	s.ableToRun = true
	return nil
}

func (s *Supervisor) Run() error {
	if s.running {
		return errors.New("CPU is already running")
	}
	if !s.ableToRun {
		return errors.New("program is not loaded or finished")
	}

	s.running = true
	for halted := false; !halted; {
		halted = s.cpu.Tick()
	}
	s.running = false
	if s.cpu.OUTP.INTN != consts.SIGTRACE {
		s.ableToRun = false
	}
	return nil
}

func (s *Supervisor) RunAll() error {
	if s.running {
		return errors.New("CPU is already running")
	}
	if !s.ableToRun {
		return errors.New("program is not loaded or finished")
	}

	s.running = true
	for finished := false; !finished; {
		for halted := false; !halted; {
			halted = s.cpu.Tick()
		}
		if s.cpu.OUTP.INTN != consts.SIGTRACE {
			finished = true
		}
	}
	s.running = false
	s.ableToRun = false
	return nil
}

func (s *Supervisor) Reset() {
	if s.running {
		s.Terminate()
	}

	s.cpu = cpu.InitCPU()
	s.compiler = compiler.Compiler{}
	s.ableToRun = false
	s.setFlags()
}

func (s *Supervisor) Terminate() error {
	s.cpu.OUTP.TERM = true
	if !s.running {
		return s.Run()
	}

	s.ableToRun = false
	return nil
}

func (s *Supervisor) TraceOn() {
	s.flags.trace = true
	s.setFlags()
}

func (s *Supervisor) SudoOn() {
	s.flags.sudo = true
	s.setFlags()
}

func (s *Supervisor) IntrOn() {
	s.flags.intr = true
	s.setFlags()
}

func (s *Supervisor) TraceOff() {
	s.flags.trace = false
	s.setFlags()
}

func (s *Supervisor) SudoOff() {
	s.flags.sudo = false
	s.setFlags()
}

func (s *Supervisor) IntrOff() {
	s.flags.intr = false
	s.setFlags()
}

func (s *Supervisor) setFlags() {
	switch s.flags.trace {
	case true:
		s.cpu.RRAM.SYS.FLG.FTOn()
	case false:
		s.cpu.RRAM.SYS.FLG.FTOff()
	}

	switch s.flags.sudo {
	case true:
		s.cpu.RRAM.SYS.FLG.FUOn()
	case false:
		s.cpu.RRAM.SYS.FLG.FUOff()
	}

	switch s.flags.intr {
	case true:
		s.cpu.RRAM.SYS.FLG.FIOn()
	case false:
		s.cpu.RRAM.SYS.FLG.FIOff()
	}
}
