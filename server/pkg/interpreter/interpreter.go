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
	AddressType          // [rx]
	ValueSourceType      // RegType OR IntType OR AddressType
	ValueDestinationType // RegType OR AddressType
)

const (
	RegTypeSize              = 6
	FlagTypeSize             = 3
	IntTypeSize              = 16
	AddressTypeSize          = 6
	ValueSourceTypeSize      = 18
	ValueDestinationTypeSize = 7
)

func SizeOfParamType(pt ParamType) int {
	switch pt {
	case RegType:
		return RegTypeSize
	case FlagType:
		return FlagTypeSize
	case IntType:
		return IntTypeSize
	case AddressType:
		return AddressTypeSize
	case ValueSourceType:
		return ValueSourceTypeSize
	case ValueDestinationType:
		return ValueDestinationTypeSize
	default:
		return 0
	}
}

func isNumType(param string) bool {
	_, err := strconv.ParseInt(param, 0, IntTypeSize+1)
	return err == nil
}

func isRegisterType(param string) bool {
	_, found := regMap[param]
	return found
}

func isAddressType(param string) bool {
	var runes = []rune(param)
	return len(param) > 1 &&
		runes[0] == '[' &&
		runes[len(runes)-1] == ']' &&
		isRegisterType(strings.Trim(param, "[]"))
}

type commandEntry struct {
	Code   types.Word32
	Params []ParamType
}

var commandsMap = map[string]commandEntry{
	"skip": {Code: types.C_SKIP, Params: []ParamType{}},
	"mov":  {Code: types.C_MOV, Params: []ParamType{ValueDestinationType, ValueSourceType}}, // dest=src
	"add":  {Code: types.C_ADD, Params: []ParamType{RegType, ValueSourceType}},              // dest+=src
	"adc":  {Code: types.C_ADC, Params: []ParamType{RegType, ValueSourceType}},              // dest+=src+fc
	"sub":  {Code: types.C_SUB, Params: []ParamType{RegType, ValueSourceType}},              // dest-=src
	"sbb":  {Code: types.C_SBB, Params: []ParamType{RegType, ValueSourceType}},              // dest-=(src+fc)
	"mul":  {Code: types.C_MUL, Params: []ParamType{RegType, ValueSourceType}},              // dest*=src
	"imul": {Code: types.C_IMUL, Params: []ParamType{RegType, ValueSourceType}},             // dest*=src (signed)
	"div":  {Code: types.C_DIV, Params: []ParamType{RegType, ValueSourceType}},              // dest/=src
	"shl":  {Code: types.C_SHL, Params: []ParamType{RegType, IntType}},                      // r<<=imm
	"shr":  {Code: types.C_SHR, Params: []ParamType{RegType, IntType}},                      // r>>=imm
	"sar":  {Code: types.C_SAR, Params: []ParamType{RegType, IntType}},                      // r<<=imm (arithmetic)
	"and":  {Code: types.C_AND, Params: []ParamType{RegType, RegType}},                      // ra&=rb
	"or":   {Code: types.C_OR, Params: []ParamType{RegType, RegType}},                       // ra|=rb
	"xor":  {Code: types.C_XOR, Params: []ParamType{RegType, RegType}},                      // ra^=rb
	"not":  {Code: types.C_NOT, Params: []ParamType{RegType}},                               // ra=~ra
	"jmp":  {Code: types.C_JMP, Params: []ParamType{AddressType}},
	/* TODO */
	"call": {Code: types.C_CALL, Params: []ParamType{}},
	"ret":  {Code: types.C_RET, Params: []ParamType{}},
	"halt": {Code: types.C_HALT, Params: []ParamType{}},
	"ei":   {Code: types.C_EI, Params: []ParamType{}},
	"di":   {Code: types.C_DI, Params: []ParamType{}},
	/*     */
	"int":  {Code: types.C_INT, Params: []ParamType{IntType}},
	"addf": {Code: types.C_ADDF, Params: []ParamType{RegType, ValueSourceType}},
	"subf": {Code: types.C_SUBF, Params: []ParamType{RegType, ValueSourceType}},
	"mulf": {Code: types.C_MULF, Params: []ParamType{RegType, ValueSourceType}},
	"divf": {Code: types.C_DIVF, Params: []ParamType{RegType, ValueSourceType}},
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

func convertRegParam(param string) (types.Word32, error) {
	rx, found := regMap[param]
	if !found {
		return 0b0, errors.New(fmt.Sprintf("register \"%s\" not found or inaccessible", param))
	}

	return rx, nil
}

func convertFlagParam(param string) (types.Word32, error) {
	fx, found := flagMap[param]
	if !found {
		return 0b0, errors.New(fmt.Sprintf("flag \"%s\" not found or inaccessible", param))
	}

	return fx, nil
}

func convertIntParam(param string) (types.Word32, error) {
	num, err := strconv.ParseInt(param, 0, IntTypeSize+1)
	if err != nil {
		return 0b0, err
	}

	return types.Word32(num), nil
}

func convertAddressParam(param string) (types.Word32, error) {
	return convertRegParam(strings.Trim(param, "[]"))
}

func convertValueSourceType(param string) (types.Word32, error) {
	switch {
	case isRegisterType(param):
		return convertRegParam(param)
	case isNumType(param):
		num, err := convertIntParam(param)
		num |= 0b01 << IntTypeSize
		return num, err
	case isAddressType(param):
		num, err := convertAddressParam(param)
		num |= 0b10 << IntTypeSize
		return num, err
	default:
		return 0b0, errors.New(fmt.Sprintf("could not convert %s to ValueSourceType", param))
	}
}

func convertParam(param string, paramType ParamType) (types.Word32, error) {
	switch paramType {
	case RegType:
		return convertRegParam(param)

	case FlagType:
		return convertFlagParam(param)

	case IntType:
		return convertIntParam(param)

	case ValueSourceType:
		return convertValueSourceType(param)

	case ValueDestinationType:
		var paramRunes = []rune(param)
		switch {
		case paramRunes[0] == 'r':
			return convertRegParam(param)

		case len(paramRunes) > 1 && paramRunes[0] == '[' && paramRunes[len(paramRunes)-1] == ']':
			sParam := string(paramRunes[1 : len(paramRunes)-1])
			num, err := convertRegParam(sParam)
			num |= 0b1 << RegTypeSize
			return num, err

		default:
			return 0b0, errors.New(fmt.Sprintf("failed to parse param %s into ValueDestinationType", param))
		}
	}
	return 0b0, errors.New(fmt.Sprintf("unsupported parameter %s with code %d", param, paramType))
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
	var shift = 5
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
