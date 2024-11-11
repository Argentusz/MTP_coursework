package consts

import "github.com/Argentusz/MTP_coursework/pkg/types"

const (
	MAX_WORD8  = 0xFF
	MAX_WORD16 = 0xFFFF
	MAX_WORD32 = 0xFFFFFFFF
)

const (
	BiGB = 1024 * 1024 * 1024 // Bytes in Gigabyte
	BiMB = 1024 * 1024
	BiKB = 1024
)

// External Memory Segments
const (
	EXE_SEG types.SegmentID = iota
	INT_SEG
	USR_SEG
	LBL_SEG
	CLL_SEG
)

const (
	SIGNONE  byte = iota
	SIGFPE        // Erroneous arithmetic operation: IGNORE
	SIGTRACE      // Breakpoint exception: WAIT
	SIGSEGV       // Invalid memory access (Segmentation fault/Stack Overflow): CRUSH
	SIGTERM       // Termination request: CRUSH
	SIGINT        // External interrupt: CRUSH
	SIGIIE        // Invalid Instruction: CRUSH
	SIGILL        // Invalid program: CRUSH
)
