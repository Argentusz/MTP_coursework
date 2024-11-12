package compiler

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"strings"
)

type Compiler struct {
	Input    []string
	Output   []types.Word32
	Aliases  map[string]string
	Labels   map[int64]types.Word32
	ILabels  map[int64]types.Word32
	Warnings []error
}

func Compile(program []string) (Compiler, error) {
	var err error

	compiler := Compiler{
		Input:   program,
		Output:  nil,
		Labels:  map[int64]types.Word32{},
		ILabels: map[int64]types.Word32{},
	}

	err = compiler.aliasStage()
	if err != nil {
		return Compiler{}, err
	}

	err = compiler.prepareStage()
	if err != nil {
		return Compiler{}, err
	}

	err = compiler.analyzeStage()
	if err != nil {
		return Compiler{}, err
	}

	err = compiler.macroStage()
	if err != nil {
		return Compiler{}, err
	}

	err = compiler.presetStage()
	if err != nil {
		return Compiler{}, err
	}

	err = compiler.compileStage()
	if err != nil {
		return Compiler{}, err
	}

	return compiler, nil
}

func (cmp *Compiler) prepareStage() error {
	program := make([]string, 0, len(cmp.Input))
	for _, v := range cmp.Input {
		line := prepLine(v)
		if line != "" {
			program = append(program, line)
		}
	}

	cmp.Input = program
	return nil
}

func (cmp *Compiler) aliasStage() error {
	cmp.setDefaultAliases()
	for i, v := range cmp.Input {
		if strings.Contains(v, "$") {
			line, err := cmp.deAlias(v)
			if err != nil {
				return err
			}

			cmp.Input[i] = line
		}
	}
	return nil
}

func (cmp *Compiler) analyzeStage() error {
	if len(cmp.Input) == 0 || cmp.Input[len(cmp.Input)-1] != "hlt" {
		cmp.Input = append(cmp.Input, "hlt")
	}

	return nil
}

func (cmp *Compiler) macroStage() error {
	program := make([]string, 0, int(1.5*float64(len(cmp.Input))))
	for _, v := range cmp.Input {
		switch isMacro(v) {
		case false:
			program = append(program, v)
		case true:
			demacro, err := deMacro(v)
			if err != nil {
				return err
			}

			program = append(program, demacro...)
		}
	}

	cmp.Input = program
	return nil
}

func (cmp *Compiler) presetStage() error {
	program := make([]string, 0, len(cmp.Input))
	for _, v := range cmp.Input {
		switch isPresetter(v) {
		case false:
			program = append(program, v)
		case true:
			err := cmp.preset((len(program))*4, v)
			if err != nil {
				return err
			}
		}
	}

	cmp.Input = program
	return nil
}

func (cmp *Compiler) compileStage() error {
	cmp.Output = make([]types.Word32, 0, len(cmp.Input))
	for _, v := range cmp.Input {
		instr, err := Convert(v)
		if err != nil {
			return err
		}

		cmp.Output = append(cmp.Output, instr)
	}

	return nil
}

func prepLine(line string) string {
	commentIdx := strings.Index(line, ";")
	if commentIdx != -1 {
		line = line[:commentIdx]
	}

	line = strings.Trim(line, " \t\n\r")
	for strings.Contains(line, "  ") {
		line = strings.ReplaceAll(line, "  ", " ")
	}

	return strings.ToLower(line)
}
