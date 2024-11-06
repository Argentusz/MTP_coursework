package register

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

type Register *types.Word32
type ShortRegister *types.Word16
type SmallRegister *types.Word8

const (
	outOfRange      byte = iota
	SmallGPRType         // 8-bit General Purpose Registers
	GPRType              // 32-bit General Purpose Registers
	ExtendedGPRType      // Extended General Purpose Registers
	HighSubGPRType       // High General Purpose Sub Registers
	LowSubGPRType        // Low General Purpose Sub Registers
)

// Registers amount by type (Subject of change)
const (
	SGPRC = 8
	GPRC  = 32
	XGPRC = 8 // Same for High/Low Sub Registers (24 Total)
)

// SysRegisters is a set of registers that are not accessible directly for user
type SysRegisters struct {
	IR   Register      // Instruction Register        - Holds memory address of current instruction
	IRB  Register      // Instruction Register Backup - For backuping IR before interrupt handling
	MBR  Register      // Memory Buffer Register      - For buffering data from External Memory
	TMP  Register      // Reserved for internal usage
	FLG  FlagsRegister // 8-bit flag register (see FlagsRegister)
	FLGI FlagsRegister // Flag register backup
}

type RRAM struct {
	SYS   SysRegisters
	SGPRs [SGPRC]SmallRegister // rb0-rb7  (IDs: 0..7)
	GPRs  [GPRC]Register       // rw0-rw31 (IDs: 8..39)
	XGPRs [XGPRC]Register      // rx0-rx7  (IDs: 40..47)
	HGPRs [XGPRC]ShortRegister // rh0-rh7  (IDs: 48..55)
	LGPRs [XGPRC]ShortRegister // rl0-rl7  (IDs: 56..63)
}

func InitRRAM() RRAM {
	var rram RRAM

	var ir, irb, mbr, tmp types.Word32

	rram.SYS.IR = &ir
	rram.SYS.IRB = &irb
	rram.SYS.MBR = &mbr
	rram.SYS.TMP = &tmp
	rram.SYS.FLG = 0b0
	rram.SYS.FLGI = 0b0

	for i := 0; i < SGPRC; i++ {
		var data types.Word8
		rram.SGPRs[i] = &data
	}

	for i := 0; i < GPRC; i++ {
		var data types.Word32
		rram.GPRs[i] = &data
	}

	for i := 0; i < XGPRC; i++ {
		XReg := InitXRegister()
		rram.XGPRs[i] = XReg.ext
		rram.HGPRs[i] = XReg.hig
		rram.LGPRs[i] = XReg.low
	}

	return rram
}

func getTypeOffset(ID types.Word32) (byte, types.Word32) {
	switch {
	case ID < SGPRC:
		return SmallGPRType, ID
	case ID < SGPRC+GPRC:
		return GPRType, ID - SGPRC
	case ID < SGPRC+GPRC+XGPRC:
		return ExtendedGPRType, ID - SGPRC - GPRC
	case ID < SGPRC+GPRC+XGPRC+XGPRC:
		return HighSubGPRType, ID - SGPRC - GPRC - XGPRC
	case ID < SGPRC+GPRC+XGPRC+XGPRC+XGPRC:
		return LowSubGPRType, ID - SGPRC - GPRC - XGPRC - XGPRC
	default:
		return outOfRange, 0
	}
}

func (rram *RRAM) GetValue(ID types.Word32) (types.Value, error) {
	regType, offset := getTypeOffset(ID)
	switch regType {
	case SmallGPRType:
		return types.Value(*rram.SGPRs[offset]), nil
	case GPRType:
		return types.Value(*rram.GPRs[offset]), nil
	case ExtendedGPRType:
		return types.Value(*rram.XGPRs[offset]), nil
	case HighSubGPRType:
		return types.Value(*rram.HGPRs[offset]), nil
	case LowSubGPRType:
		return types.Value(*rram.LGPRs[offset]), nil
	default:
		return 0, errors.New(fmt.Sprintf("unknown register ID %d", ID))
	}
}

func (rram *RRAM) PutValue(ID types.Word32, value types.Value) bool {
	regType, offset := getTypeOffset(ID)
	switch regType {
	case SmallGPRType:
		*rram.SGPRs[offset] = types.Word8(value)
		return value > consts.MAX_WORD8
	case GPRType:
		*rram.GPRs[offset] = types.Word32(value)
		return value > consts.MAX_WORD32
	case ExtendedGPRType:
		*rram.XGPRs[offset] = types.Word32(value)
		return value > consts.MAX_WORD32
	case HighSubGPRType:
		*rram.HGPRs[offset] = types.Word16(value)
		return value > consts.MAX_WORD16
	case LowSubGPRType:
		*rram.LGPRs[offset] = types.Word16(value)
		return value > consts.MAX_WORD16
	}
	return false
}

func (rram *RRAM) GetRegSize(ID types.Word32) byte {
	t, _ := getTypeOffset(ID)
	switch t {
	case SmallGPRType:
		return 8
	case GPRType:
		return 32
	case ExtendedGPRType:
		return 32
	case HighSubGPRType:
		return 16
	case LowSubGPRType:
		return 16
	default:
		return 0
	}
}
