package compiler

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/cpu"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"strings"
)

func Compile(program []string, cpu *cpu.CPU) ([]types.Word32, error) {
	demacroed := make([]string, 0, len(program))
	for _, line := range program {
		line = prepLine(line)
		if line == "" {
			continue
		}

		if isMacro(line) {
			lines, err := deMacro(line)
			if err != nil {
				return nil, err
			}

			demacroed = append(demacroed, lines...)
			continue
		}

		demacroed = append(demacroed, line)
	}
	fmt.Println(demacroed)

	compiled := make([]types.Word32, 0, len(demacroed))
	for _, line := range demacroed {
		converted, err := Convert(line)
		if err != nil {
			return nil, err
		}

		compiled = append(compiled, converted)
	}

	return compiled, nil
}

func prepLine(line string) string {
	commentIdx := strings.Index(line, ";")
	if commentIdx != -1 {
		line = line[:commentIdx]
	}

	line = strings.Trim(line, " \t\n")
	for strings.Contains(line, "  ") {
		line = strings.ReplaceAll(line, "  ", " ")
	}
	return line
}

func isMacro(line string) bool {
	return len(line) > 0 && strings.HasPrefix(line, "#")
}

func deMacro(line string) ([]string, error) {
	return []string{}, errors.New("macros are NOI")
}
