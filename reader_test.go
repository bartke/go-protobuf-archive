package pba

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestReadHeader(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	w := NewWriter(bw)

	w.WriteHeader()
	w.Flush()
	bw.Flush()

	br := bufio.NewReader(&b)
	r := NewReader(br)

	if r.Header.EntryHeaderLength != defaultEntryHeaderLength {
		t.Errorf("wrong default entry header length, expected: %v, got %v", defaultEntryHeaderLength, r.Header.EntryHeaderLength)
	}
}

func TestRead(t *testing.T) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	w := NewWriter(bw)
	w.WriteHeader()

	testData := "This is a test payload"
	w.Write([]byte(testData))

	w.Flush()
	bw.Flush()

	br := bufio.NewReader(&b)
	r := NewReader(br)

	h, err := r.Next()
	if err != nil {
		if err == io.EOF {
			t.Errorf("EOF too early")
		}
		t.Error(err)
	}

	buf := make([]byte, h.Length)
	if h.Length != uint32(len(testData)) {
		t.Errorf("payload length wrong, expected: %v, got %v", len(testData), h.Length)
	}

	_, err = r.Read(buf)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(buf, []byte(testData)) != 0 {
		t.Errorf("payload wrong, expected: %v, got %v", testData, string(buf))
	}

	_, err = r.Next()
	if err != io.EOF {
		t.Errorf("EOF expected")
	}
}
