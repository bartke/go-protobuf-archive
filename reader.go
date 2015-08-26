package pba

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	errHeader  = errors.New("pba: invalid pba header")
	errVersion = errors.New("pba: unsupported version")
)

// A Reader provides sequential access to the contents of a pba archive.
// A pba archive consists of a sequence of binary payloads.
// The Next method advances to the next payload in the archive (including the first),
// and then it can be treated as an io.Reader to access the payloads's data.
type Reader struct {
	r      io.Reader
	err    error
	curr   int // reader for current file entry
	Header *Header
}

// NewReader creates a new Reader reading from r.
func NewReader(r io.Reader) *Reader {
	tr := &Reader{
		r: r,
	}
	if tr.Header == nil {
		tr.err = tr.readHeader()
		if tr.err != nil {
			return nil
		}
	}
	return tr
}

func (tr *Reader) readHeader() error {
	tr.Header = &Header{}

	tr.err = binary.Read(tr.r, binary.LittleEndian, &tr.Header.magicNumber)
	if tr.err != nil {
		return tr.err
	}
	if bytes.Compare(tr.Header.magicNumber[:], magicNumber) != 0 {
		tr.err = errHeader
		return tr.err
	}

	tr.err = binary.Read(tr.r, binary.LittleEndian, &tr.Header.version)
	if tr.err != nil {
		return tr.err
	}
	if tr.Header.version > version {
		tr.err = errVersion
		return tr.err
	}

	tr.err = binary.Read(tr.r, binary.LittleEndian, &tr.Header.EntryHeaderLength)
	if tr.err != nil {
		return tr.err
	}

	tr.err = binary.Read(tr.r, binary.LittleEndian, &tr.Header.dataOffset)
	if tr.err != nil {
		return tr.err
	}

	tr.err = binary.Read(tr.r, binary.LittleEndian, &tr.Header.NumberOfEntries)
	if tr.err != nil {
		return tr.err
	}
	return tr.err
}

// Next advances to the next entry in the pba archive.
func (tr *Reader) Next() (*EntryHeader, error) {
	if tr.Header == nil {
		tr.err = tr.readHeader()
		if tr.err != nil {
			return nil, tr.err
		}
	}

	length := int(tr.Header.EntryHeaderLength)
	buffer := make([]byte, length)
	tr.err = binary.Read(tr.r, binary.LittleEndian, &buffer)
	if tr.err != nil {
		return nil, tr.err
	}

	eh := &EntryHeader{
		Length: littleEndian.uintBytes(buffer, length),
	}

	return eh, tr.err
}

// Read reads from the current entry in the pba archive.
// It returns 0, io.EOF when it reaches the end of that entry,
// until Next is called to advance to the next entry.
func (tr *Reader) Read(b []byte) (n int, err error) {
	n, tr.err = tr.r.Read(b)
	if tr.err != nil {
		return 0, tr.err
	}
	return n, tr.err
}
