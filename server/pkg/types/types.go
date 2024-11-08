package types

type Word8 uint8
type Word16 uint16
type Word32 uint32
type WordX interface{ Word8 | Word16 | Word32 }

type Value uint64
type SValue int64

type Address uint32
type SegmentID uint8

type ParamType byte

const (
	RegType ParamType = iota
	FlagType
	IntType
	AddressType          // [rx]
	ValueSourceType      // RegType OR IntType OR AddressType
	ValueDestinationType // RegType OR AddressType
	JumpType             // AddressType OR IntType
	LabelDestinationType // RegType OR IntType
)

const OperatorSize = 6
const (
	RegTypeSize              = 6
	FlagTypeSize             = 3
	IntTypeSize              = 16
	AddressTypeSize          = 6
	ValueSourceTypeSize      = 18
	ValueDestinationTypeSize = 7
	JumpTypeSize             = 17
	LabelDestinationTypeSize = 17
)

const (
	SourceRegMode     = 0b00 << IntTypeSize
	SourceIntMode     = 0b01 << IntTypeSize
	SourceAddrMode    = 0b10 << IntTypeSize
	SourceModeMask    = 0b11 << IntTypeSize
	SourceInverseMask = (0b1 << IntTypeSize) - 1
)

const (
	DestinationRegMode     = 0b0 << RegTypeSize
	DestinationAddrMode    = 0b1 << RegTypeSize
	DestinationModeMask    = 0b1 << RegTypeSize
	DestinationInverseMask = (0b1 << RegTypeSize) - 1
)

const (
	JumpLabelMode   = 0b0 << IntTypeSize
	JumpAddressMode = 0b1 << IntTypeSize
	JumpModeMask    = 0b1 << IntTypeSize
	JumpInverseMask = (0b1 << IntTypeSize) - 1
)

const (
	LabelIntMode     = 0b0 << IntTypeSize
	LabelRegMode     = 0b1 << IntTypeSize
	LabelModeMask    = 0b1 << IntTypeSize
	LabelInverseMask = (0b1 << IntTypeSize) - 1
)

func SizeOfParamType(pt ParamType) int {
	switch pt {
	case RegType:
		return RegTypeSize
	case FlagType:
		return FlagTypeSize
	case IntType:
		return IntTypeSize
	case AddressType:
		return AddressTypeSize
	case ValueSourceType:
		return ValueSourceTypeSize
	case ValueDestinationType:
		return ValueDestinationTypeSize
	default:
		return 0
	}
}
