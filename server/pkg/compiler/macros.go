package compiler

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"strconv"
	"strings"
	"unsafe"
)

const (
	BREAK = "#break"
	MOV32 = "#mov32"
	START = "#start"
)

func isMacro(line string) bool {
	return len(line) > 0 && strings.HasPrefix(line, "#")
}

func deMacro(line string) ([]string, error) {
	cmd := strings.Split(line, " ")
	macro := cmd[0]
	switch macro {
	case BREAK:
		return []string{"int 2"}, nil
	case START:
		return startMacro(cmd)
	case MOV32:
		return mov32Macro(cmd)
	default:
		return []string{}, errors.New(fmt.Sprintf("macro \"%s\" is NOI", macro))
	}
}

func startMacro(cmd []string) ([]string, error) {
	if len(cmd) != 2 {
		return nil, errors.New("macro #start expected label")
	}

	label := cmd[1]
	return []string{fmt.Sprintf("jmp %s", label)}, nil
}

func mov32Macro(cmd []string) ([]string, error) {
	if len(cmd) != 3 {
		return nil, errors.New("macro #mov32 expected 32-bit register and integer")
	}

	reg, snum := cmd[1], cmd[2]

	if strings.HasPrefix(snum, "f") {
		f64num, err := strconv.ParseFloat(snum[1:], 32+1)
		if err != nil {
			return nil, err
		}

		fnum := float32(f64num)
		unum := *((*uint32)(unsafe.Pointer(&fnum)))
		snum = fmt.Sprintf("0b%b", unum)
	}

	inum, err := strconv.ParseInt(snum, 0, 32+1)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(reg, "rx"):
		rln := fmt.Sprintf("rl%s", reg[2:])
		rhn := fmt.Sprintf("rh%s", reg[2:])

		rlv := inum & consts.MAX_WORD16
		rhv := (inum >> 16) & consts.MAX_WORD16

		return []string{
			fmt.Sprintf("mov %s 0b%b", rln, rlv),
			fmt.Sprintf("mov %s 0b%b", rhn, rhv),
		}, nil

	case strings.HasPrefix(reg, "rw"):
		rlv := inum & consts.MAX_WORD16
		rhv := (inum >> 16) & consts.MAX_WORD16

		return []string{
			fmt.Sprintf("mov %s 0b%b", reg, rhv),
			fmt.Sprintf("shl %s 16", reg),
			fmt.Sprintf("add %s 0b%b", reg, rlv),
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("\"%s\" is not a 32-bit register", reg))
	}
}
