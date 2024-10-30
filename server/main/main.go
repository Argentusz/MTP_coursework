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
		"mov r63 0xFF",
		"mov r63 r1",
		"mov [r63] r30",
		"mov r62 [r63]",
		"mov r62 0xFFFF",
	}
	for i := 0; i < len(lines); i++ {
		check(lines[i])
	}
}
