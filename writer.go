package pba

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	errWriteAfterClose = errors.New("pba: write after close")
	errNotImplemented  = errors.New("pba: not implemented yet")
)

// A Writer provides sequential writing of a pba archive.
// A pba archive consists of a sequence of binary payloads.
// Call WriteHeader to begin a new file, and then call Write to supply that file's data.
type Writer struct {
	w      io.Writer
	err    error
	closed bool
	Header *Header
}

// NewWriter creates a new Writer writing to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:      w,
		Header: newHeader(),
	}
}

// Flush finishes writing the current file (optional).
func (tw *Writer) Flush() error {
	return errNotImplemented
}

// WriteHeader writes hdr and prepares to accept the file's contents.
// WriteHeader calls Flush if it is not the first header.
// Calling after a Close will return ErrWriteAfterClose.
func (tw *Writer) WriteHeader() error {
	if tw.closed {
		tw.err = errWriteAfterClose
		return tw.err
	}
	if tw.err != nil {
		return tw.err
	}

	version := make([]byte, 1)
	version[0] = tw.Header.version

	entryHeaderLength := make([]byte, 1)
	entryHeaderLength[0] = tw.Header.EntryHeaderLength

	dataOffset := make([]byte, 2)
	binary.LittleEndian.PutUint16(dataOffset, tw.Header.dataOffset)

	numberOfEntries := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberOfEntries, tw.Header.NumberOfEntries)

	header := make([]byte, headerLength)
	copy(header[0:4], tw.Header.magicNumber[:])
	copy(header[4:5], version)
	copy(header[5:6], entryHeaderLength)
	copy(header[6:8], dataOffset)
	copy(header[8:16], numberOfEntries)

	_, tw.err = tw.w.Write(header)
	return tw.err
}

// Write writes to the current entry in the pba archive.
func (tw *Writer) Write(b []byte) (n int, err error) {
	if tw.closed {
		tw.err = errWriteAfterClose
		return 0, tw.err
	}
	if tw.err != nil {
		return 0, tw.err
	}

	length := int(tw.Header.EntryHeaderLength)
	header := make([]byte, length)
	littleEndian.putUintBytes(header, length, uint32(len(b)))

	_, tw.err = tw.w.Write(header)
	if tw.err != nil {
		err = tw.err
		return
	}

	n, tw.err = tw.w.Write(b)
	return n, tw.err
}

// Close closes the pba archive, flushing any unwritten
// data to the underlying writer.
func (tw *Writer) Close() error {
	if tw.err != nil || tw.closed {
		return tw.err
	}
	//tw.Flush()
	tw.closed = true
	return tw.err
}
