package main

import (
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"testing"
)

func TestRandomStuff(t *testing.T) {
	var num = 0b01111110001111000001100011111111
	fmt.Printf("%032b\n", num)

	fmt.Printf("%08b %08b %08b %08b", (num>>24)&consts.MAX_WORD8, (num>>16)&consts.MAX_WORD8, (num>>8)&consts.MAX_WORD8, (num)&consts.MAX_WORD8)
}
