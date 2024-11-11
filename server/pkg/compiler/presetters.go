package compiler

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"strconv"
	"strings"
)

const (
	LABEL  = "@label"
	ILABEL = "@ilabel"
)

func isPresetter(line string) bool {
	return len(line) > 0 && strings.HasPrefix(line, "@")
}

func (cmp *Compiler) preset(i int, v string) error {
	cmd := strings.Split(v, " ")
	if len(cmd) != 2 {
		return errors.New("presetters expect label parameter")
	}

	switch cmd[0] {
	case LABEL:
		label, err := strconv.ParseInt(cmd[1], 0, types.IntTypeSize+1)
		if err != nil {
			return err
		}

		cmp.Labels[label] = types.Word32(i)
	case ILABEL:
		ilabel, err := strconv.ParseInt(cmd[1], 0, 8+1)
		if err != nil {
			return err
		}

		cmp.ILabels[ilabel] = types.Word32(i)
	default:
		return errors.New(fmt.Sprintf("unknown presetter \"%s\"", cmd[0]))
	}
	return nil
}
