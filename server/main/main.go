package main

import (
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/interpreter"
	"strconv"
)

func check(line string) {
	cmd, err := interpreter.Convert(line)
	if err != nil {
		fmt.Println("Error while parsing \"", line, "\":", err.Error())
		return
	}
	fmt.Printf("%s: %034s\n", line, strconv.FormatInt(int64(cmd), 2))
}

func main() {
	var lines = []string{
		"skp", "inc r1", "ass r2 2",
		"skp", "inc r63", "ass r3 0b101",
		"skp 1", "inc r65", "ass r4 0b11111111111111111111111111111111111",
		"ass fo 1", "inc r1 r2", "ass", "inc r64",
	}
	for i := 0; i < len(lines); i++ {
		check(lines[i])
	}
}
