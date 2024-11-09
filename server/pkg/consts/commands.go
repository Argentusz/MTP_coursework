package consts

import "github.com/Argentusz/MTP_coursework/pkg/types"

// Instruction codes
const (
	C_SKIP types.Word32 = iota
	C_MOV
	C_ADD
	C_ADC
	C_SUB
	C_SBB
	C_MUL
	C_DIV
	C_IMOV
	C_IADD
	C_IADC
	C_ISUB
	C_ISBB
	C_IMUL
	C_IDIV
	C_ADDF
	C_SUBF
	C_MULF
	C_DIVF
	C_SHL
	C_SHR
	C_SAR
	C_AND
	C_OR
	C_NOT
	C_XOR
	C_JMP
	C_JIF
	C_JNF
	C_JIZ
	C_JNZ
	C_LBL
	C_CALL
	C_RET
	C_EI
	C_DI
	C_INT
	C_HLT = 0o77
)
