package pba

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"
)

func TestWriteHeader(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	w := NewWriter(bw)
	w.Header.NumberOfEntries = 1234567890

	w.WriteHeader()
	w.Flush()

	bw.Flush()

	// check magic number
	mn := make([]byte, 4)
	b.Read(mn)
	if bytes.Compare(mn, magicNumber) != 0 {
		t.Errorf("wrong magic number, expected: %v, got %v", magicNumber, mn)
	}

	// check version
	v := make([]byte, 1)
	b.Read(v)
	if bytes.Compare(v, []byte{version}) != 0 {
		t.Errorf("wrong version, expected: %v, got %v", version, v)
	}

	// check payload header length
	l := make([]byte, 1)
	b.Read(l)
	if bytes.Compare(l, []byte{w.Header.EntryHeaderLength}) != 0 {
		t.Errorf("wrong payload header length, expected: %v, got %v", w.Header.EntryHeaderLength, l)
	}

	// check data offset
	offset := make([]byte, 2)
	b.Read(offset)
	offsetN := binary.LittleEndian.Uint16(offset)
	if offsetN != headerLength {
		t.Errorf("wrong data offset, expected: %v, got %v", headerLength, offsetN)
	}

	// check number of entries
	entries := make([]byte, 8)
	b.Read(entries)
	entriesN := binary.LittleEndian.Uint64(entries)
	if offsetN != headerLength {
		t.Errorf("wrong number of entries, expected: %v, got %v", 1234567890, entriesN)
	}
}

func TestWriter(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	w := NewWriter(bw)
	w.WriteHeader()

	testData := "This is a test payload"
	w.Write([]byte(testData))

	w.Flush()
	bw.Flush()

	rawHeader := make([]byte, 16)
	b.Read(rawHeader)

	rawEntryHeader := make([]byte, w.Header.EntryHeaderLength)
	b.Read(rawEntryHeader)

	l := binary.LittleEndian.Uint16(rawEntryHeader)
	if l != uint16(len(testData)) {
		t.Errorf("payload length wrong, expected: %v, got %v", len(testData), l)
	}

	rawPayload := make([]byte, len(testData))
	b.Read(rawPayload)

	if bytes.Compare(rawPayload, []byte(testData)) != 0 {
		t.Errorf("payload wrong, expected: %v, got %v", testData, string(rawPayload))
	}
}
