package types

type Word8 uint8
type Word16 uint16
type Word32 uint32
type WordX interface{ Word8 | Word16 | Word32 }

type Address uint32
type SegmentID uint8
