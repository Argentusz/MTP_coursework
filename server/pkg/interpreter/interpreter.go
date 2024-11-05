package interpreter

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"strconv"
	"strings"
)

func isNumType(param string) bool {
	_, err := strconv.ParseInt(param, 0, types.IntTypeSize+1)
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
	Params []types.ParamType
}

var commandsMap = map[string]commandEntry{
	"skip": {Code: consts.C_SKIP, Params: []types.ParamType{}},
	"mov":  {Code: consts.C_MOV, Params: []types.ParamType{types.ValueDestinationType, types.ValueSourceType}}, // dest=src
	"add":  {Code: consts.C_ADD, Params: []types.ParamType{types.RegType, types.ValueSourceType}},              // dest+=src
	"adc":  {Code: consts.C_ADC, Params: []types.ParamType{types.RegType, types.ValueSourceType}},              // dest+=src+fc
	"sub":  {Code: consts.C_SUB, Params: []types.ParamType{types.RegType, types.ValueSourceType}},              // dest-=src
	"sbb":  {Code: consts.C_SBB, Params: []types.ParamType{types.RegType, types.ValueSourceType}},              // dest-=(src+fc)
	"mul":  {Code: consts.C_MUL, Params: []types.ParamType{types.RegType, types.ValueSourceType}},              // dest*=src
	"div":  {Code: consts.C_DIV, Params: []types.ParamType{types.RegType, types.ValueSourceType}},              // dest/=src
	"iadd": {Code: consts.C_IADD, Params: []types.ParamType{types.RegType, types.ValueSourceType}},             // dest+=src
	"iadc": {Code: consts.C_IADC, Params: []types.ParamType{types.RegType, types.ValueSourceType}},             // dest+=src+fc
	"isub": {Code: consts.C_ISUB, Params: []types.ParamType{types.RegType, types.ValueSourceType}},             // dest-=src
	"isbb": {Code: consts.C_ISBB, Params: []types.ParamType{types.RegType, types.ValueSourceType}},             // dest-=(src+fc)
	"imul": {Code: consts.C_IMUL, Params: []types.ParamType{types.RegType, types.ValueSourceType}},             // dest*=src (signed)
	"idiv": {Code: consts.C_IDIV, Params: []types.ParamType{types.RegType, types.ValueSourceType}},             // dest/=src
	"shl":  {Code: consts.C_SHL, Params: []types.ParamType{types.RegType, types.IntType}},                      // r<<=imm
	"shr":  {Code: consts.C_SHR, Params: []types.ParamType{types.RegType, types.IntType}},                      // r>>=imm
	"sar":  {Code: consts.C_SAR, Params: []types.ParamType{types.RegType, types.IntType}},                      // r<<=imm (arithmetic)
	"and":  {Code: consts.C_AND, Params: []types.ParamType{types.RegType, types.RegType}},                      // ra&=rb
	"or":   {Code: consts.C_OR, Params: []types.ParamType{types.RegType, types.RegType}},                       // ra|=rb
	"xor":  {Code: consts.C_XOR, Params: []types.ParamType{types.RegType, types.RegType}},                      // ra^=rb
	"not":  {Code: consts.C_NOT, Params: []types.ParamType{types.RegType}},                                     // ra=~ra
	"jmp":  {Code: consts.C_JMP, Params: []types.ParamType{types.AddressType}},
	/* TODO */
	"call": {Code: consts.C_CALL, Params: []types.ParamType{}},
	"ret":  {Code: consts.C_RET, Params: []types.ParamType{}},
	"halt": {Code: consts.C_HALT, Params: []types.ParamType{}},
	"ei":   {Code: consts.C_EI, Params: []types.ParamType{}},
	"di":   {Code: consts.C_DI, Params: []types.ParamType{}},
	/*     */
	"int":  {Code: consts.C_INT, Params: []types.ParamType{types.IntType}},
	"addf": {Code: consts.C_ADDF, Params: []types.ParamType{types.RegType, types.ValueSourceType}},
	"subf": {Code: consts.C_SUBF, Params: []types.ParamType{types.RegType, types.ValueSourceType}},
	"mulf": {Code: consts.C_MULF, Params: []types.ParamType{types.RegType, types.ValueSourceType}},
	"divf": {Code: consts.C_DIVF, Params: []types.ParamType{types.RegType, types.ValueSourceType}},
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
	// TODO: No Magic numbers connect to consts
	var id types.Word32
	for i := types.Word32(0); i < 8; i++ {
		regMap[fmt.Sprintf("rb%d", i)] = id
		id++
	}
	for i := types.Word32(0); i < 32; i++ {
		regMap[fmt.Sprintf("rw%d", i)] = id
		id++
	}
	for i := types.Word32(0); i < 8; i++ {
		regMap[fmt.Sprintf("rx%d", i)] = id
		id++
	}
	for i := types.Word32(0); i < 8; i++ {
		regMap[fmt.Sprintf("rh%d", i)] = id
		id++
	}
	for i := types.Word32(0); i < 8; i++ {
		regMap[fmt.Sprintf("rl%d", i)] = id
		id++
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
	num, err := strconv.ParseInt(param, 0, types.IntTypeSize+1)
	if err != nil {
		return 0b0, err
	}

	if num < 0 {
		return types.Word32(-num) | (1 << (types.IntTypeSize - 1)), nil
	}

	return types.Word32(num), nil
}

func convertAddressParam(param string) (types.Word32, error) {
	return convertRegParam(strings.Trim(param, "[]"))
}

func convertValueSourceType(param string) (types.Word32, error) {
	switch {
	case isRegisterType(param):
		num, err := convertRegParam(param)
		num |= types.SourceRegMode
		return num, err
	case isNumType(param):
		num, err := convertIntParam(param)
		num |= types.SourceIntMode
		return num, err
	case isAddressType(param):
		num, err := convertAddressParam(param)
		num |= types.SourceAddrMode
		return num, err
	default:
		return 0b0, errors.New(fmt.Sprintf("could not convert %s to ValueSourceType", param))
	}
}

func convertParam(param string, paramType types.ParamType) (types.Word32, error) {
	switch paramType {
	case types.RegType:
		return convertRegParam(param)

	case types.FlagType:
		return convertFlagParam(param)

	case types.IntType:
		return convertIntParam(param)

	case types.ValueSourceType:
		return convertValueSourceType(param)

	case types.ValueDestinationType:
		var paramRunes = []rune(param)
		switch {
		case paramRunes[0] == 'r':
			return convertRegParam(param)

		case len(paramRunes) > 1 && paramRunes[0] == '[' && paramRunes[len(paramRunes)-1] == ']':
			sParam := string(paramRunes[1 : len(paramRunes)-1])
			num, err := convertRegParam(sParam)
			num |= 0b1 << types.RegTypeSize
			return num, err

		default:
			return 0b0, errors.New(fmt.Sprintf("failed to parse param %s into ValueDestinationType", param))
		}
	}
	return 0b0, errors.New(fmt.Sprintf("unsupported parameter %s with code %d", param, paramType))
}

func Convert(instr string) (types.Word32, error) {
	//fmt.Println("Converting", instr)
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
	//fmt.Printf("Command code: %05b\n\n", cmd)
	var shift = types.OperatorSize
	for i := 0; i < len(entry.Params); i++ {
		//fmt.Println("Converting param: ", command[i+1])
		param, err := convertParam(command[i+1], entry.Params[i])
		//fmt.Printf(fmt.Sprintf("param code: %%0%db\n", types.SizeOfParamType(entry.Params[i])), param)
		if err != nil {
			return 0b0, err
		}

		cmd |= param << shift
		//fmt.Printf("Resulting in cmd: %032b\n", cmd)
		shift += types.SizeOfParamType(entry.Params[i])
		//fmt.Println()
	}

	if shift > 32 {
		return 0, errors.New("command exceeds 32-bit limit")
	}
	return cmd, nil
}
