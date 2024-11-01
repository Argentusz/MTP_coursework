package xmem

import (
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"testing"
)

const testSegment types.SegmentID = 1

func TestXMem(t *testing.T) {
	xmem := InitExternalMemory()

	err := xmem.NewSegment(testSegment, 0b111)
	if err != nil {
		t.Fatal("Error while creating new segment\n")
	}
	t.Log("Successfully initialized xmem")

	err = xmem.NewSegment(testSegment, 0b111)
	if err == nil {
		t.Error("Segment overwrite\n")
	}
	t.Log("Successfully initialized memseg")

	data, err := xmem.At(testSegment).GetByte(0b0)
	if data != 0 || err != nil {
		t.Error("Can not read uninitialized data")
		t.Error("Error:", err.Error())
	}
	t.Log("Successfully read uninitialized data")

	data32, err := xmem.At(testSegment).GetWord32(0b0)
	if data32 != 0 || err != nil {
		t.Error("Can not read uninitialized 32 data")
		t.Error("Error:", err.Error())
	}
	t.Log("Successfully read 32 uninitialized data")

	data, err = xmem.At(testSegment).GetByte(0b1000)
	if data != 0 || err == nil {
		t.Error("Read out of bounds")
	}
	t.Log("Successfully avoided out of bounds read")

	data32, err = xmem.At(testSegment).GetWord32(0b111)
	if data != 0 || err == nil {
		t.Error("Read 32 out of bounds")
	}
	t.Log("Successfully avoided out of bounds 32 read")

	var inputData types.Word8 = 0b10100101
	err = xmem.At(testSegment).SetByte(0b1, inputData)
	if err != nil {
		t.Fatal("Can not write")
	}
	data, err = xmem.At(testSegment).GetByte(0b1)
	if data != inputData || err != nil {
		t.Error("Read/Write do not match")
	}
	t.Log("Successfully wrote and read same data")

	var inputData32 types.Word32 = 0b10100101
	err = xmem.At(testSegment).SetWord32(0b1, inputData32)
	if err != nil {
		t.Fatal("Can not write 32")
	}
	data32, err = xmem.At(testSegment).GetWord32(0b1)
	if data != inputData || err != nil {
		t.Error("Read/Write 32 do not match")
	}
	t.Log("Successfully wrote and read same data")

	err = xmem.At(testSegment).SetByte(0b1000, inputData)
	if err == nil {
		t.Error("Written out of bounds")
	}
	t.Log("Successfully avoided writing out of bounds")

	err = xmem.At(testSegment).SetWord32(0b101, inputData32)
	if err == nil {
		t.Error("Written out of bounds")
	}
	t.Log("Successfully avoided writing out of bounds")
}
