package interpreter

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"strconv"
	"strings"
)

type ParamType byte

const (
	RegType ParamType = iota
	FlagType
	IntType
	FloatType
)

const (
	RegTypeSize   = 6
	FlagTypeSize  = 3
	IntTypeSize   = 16
	FloatTypeSize = 16
)

func SizeOfParamType(pt ParamType) int {
	switch pt {
	case RegType:
		return RegTypeSize
	case FlagType:
		return FlagTypeSize
	case IntType:
		return IntTypeSize
	case FloatType:
		return FloatTypeSize
	default:
		return 0
	}
}

type commandEntry struct {
	Code   types.Word32
	Params []ParamType
}

var commandsMap = map[string]commandEntry{
	"skp": {Code: 0b00000000, Params: []ParamType{}},
	"inc": {Code: 0b00000001, Params: []ParamType{RegType}},
	"ass": {Code: 0b00000010, Params: []ParamType{RegType, IntType}},
}

var flagMap = map[string]types.Word32{
	"fz": 0b000,
	"fc": 0b001,
	"fs": 0b010,
	"fo": 0b011,
	"fi": 0b100,
	"ft": 0b101,
	"fu": 0b110,
	"fx": 0b111,
}

var regMap = map[string]types.Word32{}

func init() {
	for i := types.Word32(0); i < 64; i++ {
		regMap[fmt.Sprintf("r%d", i)] = i
	}
}

func convertParam(param string, paramType ParamType) (types.Word32, error) {
	switch paramType {
	case RegType:
		rx, found := regMap[param]
		if !found {
			return 0b0, errors.New(fmt.Sprintf("register \"%s\" not found or inaccessible", param))
		}

		return rx, nil
	case IntType:
		num, err := strconv.ParseInt(param, 0, IntTypeSize)
		if err != nil {
			return 0b0, err
		}

		return types.Word32(num), nil
	}
	return 0b0, errors.New(fmt.Sprintf("unsupported parameter type with code %d", paramType))
}

func Convert(instr string) (types.Word32, error) {
	command := strings.Split(instr, " ")
	entry, found := commandsMap[command[0]]
	if !found {
		return 0b0, errors.New(fmt.Sprintf("command \"%s\" not found", command[0]))
	}

	if len(command)-1 != len(entry.Params) {
		return 0b0, errors.New(fmt.Sprintf("command \"%s\" expected %d parameters, got %d", command[0], len(entry.Params), len(command)-1))
	}

	var cmd types.Word32
	cmd = entry.Code
	var shift = 8
	for i := 0; i < len(entry.Params); i++ {
		param, err := convertParam(command[i+1], entry.Params[i])
		if err != nil {
			return 0b0, err
		}

		cmd |= param << shift
		shift += SizeOfParamType(entry.Params[i])
	}

	if shift > 32 {
		return 0, errors.New("command exceeds 32-bit limit")
	}

	return cmd, nil
}
