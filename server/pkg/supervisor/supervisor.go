package supervisor

import (
	"errors"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
)

type Supervisor struct {
	cpu   cpu.CPU
	trace bool
	sudo  bool
}

func (s *Supervisor) Initialize(trace, sudo bool) {
	s.cpu = cpu.InitCPU()
	s.trace = trace
	s.sudo = sudo
}

func (s *Supervisor) Reset() {
	s.cpu = cpu.InitCPU()
}

func (s *Supervisor) Terminate() {
	s.cpu.ForceSIGINT()
}

func (s *Supervisor) Step() error {
	if !s.trace {
		return errors.New("can not step in non-trace mode")
	}

	bounder := 0
	for halted := false; !halted; bounder++ {
		halted = s.cpu.Tick()
		if bounder > 50 {
			return errors.New("unhalted step")
		}
	}

	return nil
}

func (s *Supervisor) RunAll() error {

	return nil
}
