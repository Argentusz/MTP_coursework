package xmem

import (
	"errors"
	"github.com/Argentusz/MTP_coursework/pkg/types"
)

type ExternalMemory struct {
	Segments map[types.SegmentID]MemorySegment
}

func InitExternalMemory() ExternalMemory {
	return ExternalMemory{Segments: map[types.SegmentID]MemorySegment{}}
}

func (xmem *ExternalMemory) NewSegment(ID types.SegmentID, bytes uint32) error {
	_, found := xmem.Segments[ID]
	if found {
		return errors.New("memory with given ID already exists")
	}

	xmem.Segments[ID] = InitSegment(0b0, types.Address(bytes)-1)
	return nil
}

func (xmem *ExternalMemory) At(ID types.SegmentID) *MemorySegment {
	seg := xmem.Segments[ID]
	return &seg
}
