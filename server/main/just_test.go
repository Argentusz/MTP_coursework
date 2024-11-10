package main

import (
	"github.com/Argentusz/MTP_coursework/pkg/compiler"
	"github.com/Argentusz/MTP_coursework/pkg/types"
	"testing"
)

func TestRandomStuff(t *testing.T) {
	int_ignore_mtp := []string{
		"skip",
		"ret",
	}
	int_finish_mtp := []string{
		"hlt",
		"ret",
	}

	var int_ignore []types.Word32
	for _, line := range int_ignore_mtp {
		c, _ := compiler.Convert(line)
		int_ignore = append(int_ignore, c)
	}

	var int_finish []types.Word32
	for _, line := range int_finish_mtp {
		c, _ := compiler.Convert(line)
		int_finish = append(int_finish, c)
	}

	t.Log(int_ignore)
	t.Log(int_finish)
}
