package compiler

import "testing"

func TestPrepLine(t *testing.T) {
	type entry struct {
		line     string
		expected string
	}

	entries := []entry{
		{line: "", expected: ""},
		{line: "          ", expected: ""},
		{line: "\t", expected: ""},
		{line: "; Comment", expected: ""},
		{line: "    ;;;;;; Comment ;;;;;;", expected: ""},
		{line: "mov rb0 1", expected: "mov rb0 1"},
		{line: "     mov        rb0        1     ", expected: "mov rb0 1"},
		{line: "mov rb0 1 ; Comment", expected: "mov rb0 1"},
		{line: "mov     rb0   1  ;      Comment", expected: "mov rb0 1"},
	}

	for i, entry := range entries {
		line := entry.line
		expected := entry.expected

		if prepLine(line) != expected {
			t.Errorf("[ERROR] Line #%d: Got \"%s\"; Expected \"%s\"", i, line, expected)
		} else {
			t.Logf("Line #%d compiled successfully: \"%s\"", i, line)
		}
	}
}
