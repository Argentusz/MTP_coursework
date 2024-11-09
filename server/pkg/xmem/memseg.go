package xmem

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

type MemorySegment struct {
	Table   map[types.Address]types.Word8
	minAddr types.Address
	maxAddr types.Address
}

func InitSegment(minAddr, maxAddr types.Address) MemorySegment {
	return MemorySegment{
		Table:   map[types.Address]types.Word8{},
		minAddr: minAddr,
		maxAddr: maxAddr,
	}
}

func (mseg *MemorySegment) GetByte(addr types.Address) (types.Word8, error) {
	if addr > mseg.maxAddr || addr < mseg.minAddr {
		return 0, errors.New(fmt.Sprintf("address %o is out of range for segment", addr))
	}

	return mseg.Table[addr], nil
}

func (mseg *MemorySegment) SetByte(addr types.Address, data types.Word8) error {
	if addr > mseg.maxAddr || addr < mseg.minAddr {
		return errors.New(fmt.Sprintf("address %o is out of range for segment", addr))
	}

	mseg.Table[addr] = data
	return nil
}

func (mseg *MemorySegment) GetWord32(addr types.Address) (types.Word32, error) {
	if addr+3 > mseg.maxAddr || addr < mseg.minAddr {
		return 0, errors.New(fmt.Sprintf("address %o is out of range for segment", addr))
	}

	var word32 types.Word32
	for i := types.Address(0); i < 4; i++ {
		word32 |= types.Word32(mseg.Table[addr+i]) << (8 * (3 - i))
	}
	return word32, nil
}

func (mseg *MemorySegment) SetWord32(addr types.Address, data types.Word32) error {
	if addr+3 > mseg.maxAddr || addr < mseg.minAddr {
		return errors.New(fmt.Sprintf("address %o is out of range for segment", addr))
	}

	for i := types.Address(0); i < 4; i++ {
		word8 := types.Word8((data >> ((3 - i) * 8)) & consts.MAX_WORD8)
		mseg.Table[addr+i] = word8
	}
	return nil
}

func (mseg *MemorySegment) SetWord16(addr types.Address, data types.Word16) error {
	if addr+1 > mseg.maxAddr || addr < mseg.minAddr {
		return errors.New(fmt.Sprintf("address %o is out of range for segment", addr))
	}

	for i := types.Address(0); i < 2; i++ {
		word8 := types.Word8((data >> ((1 - i) * 8)) & consts.MAX_WORD8)
		mseg.Table[addr+i] = word8
	}
	return nil
}

func (mseg *MemorySegment) GetMaxAddr() types.Address {
	return mseg.maxAddr
}
