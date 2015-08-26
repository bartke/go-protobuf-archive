// Package pba implements access to pba archives.
package pba

import "encoding/binary"

const (
	version      = 0 // unversioned WIP
	headerLength = 16

	defaultEntryHeaderLength = 2
)

var (
	magicNumber = []byte{0x4C, 0x01, 0x2A, 0x00}
)

// A Header represents the file header in a pba archive.
type Header struct {
	magicNumber [4]byte
	version     byte
	dataOffset  uint16

	NumberOfEntries   uint64
	EntryHeaderLength byte
}

func newHeader() *Header {
	header := &Header{
		version:           version,
		dataOffset:        headerLength,
		EntryHeaderLength: defaultEntryHeaderLength,
	}

	copy(header.magicNumber[:], magicNumber)

	return header
}

// EntryHeader represents a single payload header in a pba archive.
type EntryHeader struct {
	Length uint32
}

type le struct{}

var littleEndian le

func (le) uint24(b []byte) uint32 {
	return uint32(b[2])<<16 |
		uint32(b[1])<<8 |
		uint32(b[0])
}

func (le) putUint24(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func (le) int64(b []byte) int64 {
	return int64(b[7])<<56 |
		int64(b[6])<<48 |
		int64(b[5])<<40 |
		int64(b[4])<<32 |
		int64(b[3])<<24 |
		int64(b[2])<<16 |
		int64(b[1])<<8 |
		int64(b[0])
}

func (le) putInt64(b []byte, v int64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func (le) uintBytes(b []byte, nBytes int) uint32 {
	switch nBytes {
	case 1:
		return uint32(b[0])
	case 2:
		return uint32(binary.LittleEndian.Uint16(b))
	case 3:
		return littleEndian.uint24(b)
	default:
		return binary.LittleEndian.Uint32(b)
	}
}

func (le) putUintBytes(b []byte, nBytes int, v uint32) {
	switch nBytes {
	case 1:
		b[0] = uint8(v)
	case 2:
		binary.LittleEndian.PutUint16(b, uint16(v))
	case 3:
		littleEndian.putUint24(b, v)
	default:
		binary.LittleEndian.PutUint32(b, v)
	}
}
